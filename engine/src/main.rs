use std::io::{self, Read};
use anyhow::Result;
use skene_engine::{
    protocol::{EngineInput, EngineOutput},
    codebase::CodebaseExplorer,
    llm::create_llm_client,
    analyzers::{
        tech_stack::create_tech_stack_analyzer,
        growth_features::create_growth_features_analyzer,
        revenue_leakage::create_revenue_leakage_analyzer,
        industry::create_industry_analyzer,
        manifest::create_manifest_analyzer,
        docs::{create_product_overview_analyzer, create_features_analyzer, create_docs_manifest_analyzer},
    },
    planner::Planner,
    output::{write_manifest_json, write_file, write_product_docs},
    manifest::GrowthManifest,
    strategies::context::AnalysisContext,
};
use std::path::PathBuf;
use serde_json::json;

#[tokio::main]
async fn main() {
    if let Err(e) = run().await {
        let error_out = EngineOutput::Error {
            message: e.to_string(),
            code: None,
        };
        println!("{}", serde_json::to_string(&error_out).unwrap());
        std::process::exit(1);
    }
}

async fn run() -> Result<()> {
    let stdin = io::stdin();
    let mut handle = stdin.lock();
    let mut input_str = String::new();
    handle.read_to_string(&mut input_str)?;

    let input: EngineInput = serde_json::from_str(&input_str)?;

    match input.command.as_str() {
        "analyze" => run_analyze(input).await?,
        "plan" => run_plan(input).await?,
        _ => return Err(anyhow::anyhow!("Unknown command: {}", input.command)),
    }

    Ok(())
}

async fn run_analyze(input: EngineInput) -> Result<()> {
    let base_dir = PathBuf::from(&input.project_dir);
    let output_dir = PathBuf::from(&input.output_dir);
    let explorer = CodebaseExplorer::new(base_dir.clone(), Some(input.exclude_folders));
    let llm_client = create_llm_client(&input.provider, &input.api_key, &input.model, input.base_url.as_deref())
        .map_err(|e| anyhow::anyhow!(e))?;
    let llm = llm_client.as_ref();

    let on_progress = |phase: &str, msg: String, p: f64| {
        let output = EngineOutput::Progress {
            phase: phase.to_string(),
            step: 0, 
            total_steps: 0,
            progress: p,
            message: msg,
        };
        println!("{}", serde_json::to_string(&output).unwrap());
    };

    // Run core analyzers
    on_progress("tech_stack", "Analyzing tech stack...".to_string(), 0.1);
    let ts_result = create_tech_stack_analyzer().run(&explorer, llm, "Detect tech stack".to_string(), |m, _, _, _| on_progress("tech_stack", m, 0.2)).await?;
    
    on_progress("growth_features", "Analyzing growth features...".to_string(), 0.3);
    let gf_result = create_growth_features_analyzer().run(&explorer, llm, "Detect growth features".to_string(), |m, _, _, _| on_progress("growth_features", m, 0.4)).await?;
    
    on_progress("revenue_leakage", "Analyzing revenue leakage...".to_string(), 0.5);
    let rl_result = create_revenue_leakage_analyzer().run(&explorer, llm, "Detect revenue leakage".to_string(), |m, _, _, _| on_progress("revenue_leakage", m, 0.6)).await?;
    
    on_progress("industry", "Analyzing industry...".to_string(), 0.7);
    let ind_result = create_industry_analyzer().run(&explorer, llm, "Detect industry".to_string(), |m, _, _, _| on_progress("industry", m, 0.8)).await?;

    // Prepare context for manifest
    let mut manifest_context = AnalysisContext::new("Generate Manifest".to_string());
    if let Some(ts) = ts_result.get("tech_stack") { manifest_context.set("tech_stack", ts.clone()); }
    if let Some(gf) = gf_result.get("current_growth_features") { manifest_context.set("current_growth_features", gf.clone()); }
    if let Some(rl) = rl_result.get("revenue_leakage") { manifest_context.set("revenue_leakage", rl.clone()); }
    if let Some(ind) = ind_result.get("industry") { manifest_context.set("industry", ind.clone()); }

    let manifest_path = output_dir.join("growth-manifest.json");
    let mut docs_path = None;

    let manifest_json = if input.product_docs {
        // Run doc analyzers
        on_progress("product_overview", "Analyzing product overview...".to_string(), 0.85);
        let po_result = create_product_overview_analyzer().run(&explorer, llm, "Product Overview".to_string(), |_,_,_,_| {}).await?;
        
        on_progress("features", "Documenting features...".to_string(), 0.9);
        let feat_result = create_features_analyzer().run(&explorer, llm, "Features".to_string(), |_,_,_,_| {}).await?;
        
        if let Some(po) = po_result.get("product_overview") { manifest_context.set("product_overview", po.clone()); }
        if let Some(feat) = feat_result.get("features") { manifest_context.set("features", feat.clone()); }
        
        on_progress("manifest", "Generating docs manifest...".to_string(), 0.95);
        let docs_res = create_docs_manifest_analyzer().run_with_context(&explorer, llm, "Generate Docs Manifest".to_string(), Some(manifest_context), |m, _, _, _| on_progress("manifest", m, 0.95)).await?;
        
        if let Some(json) = docs_res.get("docs_manifest") {
            let p_docs_path = output_dir.join("product-docs.md");
            write_product_docs(&p_docs_path, json).await?;
            docs_path = Some(p_docs_path.to_string_lossy().to_string());
            Some(json.clone())
        } else {
            None
        }
    } else {
        on_progress("manifest", "Generating manifest...".to_string(), 0.95);
        let man_res = create_manifest_analyzer().run_with_context(&explorer, llm, "Generate Manifest".to_string(), Some(manifest_context), |m, _, _, _| on_progress("manifest", m, 0.95)).await?;
        man_res.get("manifest").cloned()
    };

    if let Some(json) = manifest_json {
        write_manifest_json(&manifest_path, &json).await?;
        
        let output = EngineOutput::Result {
            manifest_path: Some(manifest_path.to_string_lossy().to_string()),
            template_path: None,
            docs_path,
            plan_path: None,
        };
        println!("{}", serde_json::to_string(&output).unwrap());
    } else {
        return Err(anyhow::anyhow!("Failed to generate manifest"));
    }
    
    Ok(())
}

async fn run_plan(input: EngineInput) -> Result<()> {
    // Read manifest
    let manifest_path = input.manifest_path.ok_or_else(|| anyhow::anyhow!("Manifest path required for plan"))?;
    let manifest_str = tokio::fs::read_to_string(&manifest_path).await?;
    let manifest: GrowthManifest = serde_json::from_str(&manifest_str)?;
    
    let llm_client = create_llm_client(&input.provider, &input.api_key, &input.model, input.base_url.as_deref())
        .map_err(|e| anyhow::anyhow!(e))?;
    
    let output_path = PathBuf::from(&input.output_dir).join("growth-plan.md");
    
    let plan_content = if input.onboarding.unwrap_or(false) {
        Planner::generate_onboarding_memo(llm_client.as_ref(), &manifest).await?
    } else {
        Planner::generate_council_memo(llm_client.as_ref(), &manifest).await?
    };
    
    write_file(&output_path, &plan_content).await?;
    
    let output = EngineOutput::Result {
        manifest_path: Some(manifest_path),
        template_path: None,
        docs_path: None,
        plan_path: Some(output_path.to_string_lossy().to_string()),
    };
    println!("{}", serde_json::to_string(&output).unwrap());
    
    Ok(())
}
