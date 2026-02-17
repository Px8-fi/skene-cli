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
        "build" => run_build(input).await?,
        "status" => run_status(input).await?,
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

    on_progress("plan", "Loading manifest...".to_string(), 0.1);

    // Read manifest
    let manifest_path = input.manifest_path.ok_or_else(|| anyhow::anyhow!("Manifest path required for plan"))?;
    let manifest_str = tokio::fs::read_to_string(&manifest_path).await?;
    let manifest: GrowthManifest = serde_json::from_str(&manifest_str)?;

    on_progress("plan", "Connecting to LLM provider...".to_string(), 0.2);

    let llm_client = create_llm_client(&input.provider, &input.api_key, &input.model, input.base_url.as_deref())
        .map_err(|e| anyhow::anyhow!(e))?;
    
    let output_path = PathBuf::from(&input.output_dir).join("growth-plan.md");

    on_progress("plan", "Generating growth plan...".to_string(), 0.4);
    
    let plan_content = if input.onboarding.unwrap_or(false) {
        Planner::generate_onboarding_memo(llm_client.as_ref(), &manifest).await?
    } else {
        Planner::generate_council_memo(llm_client.as_ref(), &manifest).await?
    };

    on_progress("plan", "Writing plan to disk...".to_string(), 0.9);
    
    write_file(&output_path, &plan_content).await?;

    on_progress("plan", "Growth plan generated".to_string(), 1.0);
    
    let output = EngineOutput::Result {
        manifest_path: Some(manifest_path),
        template_path: None,
        docs_path: None,
        plan_path: Some(output_path.to_string_lossy().to_string()),
    };
    println!("{}", serde_json::to_string(&output).unwrap());
    
    Ok(())
}

async fn run_build(input: EngineInput) -> Result<()> {
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

    on_progress("build", "Loading manifest...".to_string(), 0.1);

    // Read manifest
    let output_dir = PathBuf::from(&input.output_dir);
    let manifest_path = input.manifest_path
        .map(PathBuf::from)
        .unwrap_or_else(|| output_dir.join("growth-manifest.json"));

    let manifest_str = tokio::fs::read_to_string(&manifest_path).await
        .map_err(|_| anyhow::anyhow!("No manifest found at {}. Run analysis first.", manifest_path.display()))?;
    let manifest: GrowthManifest = serde_json::from_str(&manifest_str)?;

    on_progress("build", "Connecting to LLM provider...".to_string(), 0.2);

    let llm_client = create_llm_client(&input.provider, &input.api_key, &input.model, input.base_url.as_deref())
        .map_err(|e| anyhow::anyhow!(e))?;

    on_progress("build", "Generating implementation prompt...".to_string(), 0.4);

    let manifest_summary = format_manifest_for_prompt(&manifest);
    let prompt = format!(
        "You are an expert AI coding assistant prompt engineer.\n\n\
        Based on the following growth manifest, generate a detailed implementation prompt \
        that a developer can paste directly into Cursor, Claude, or another AI coding tool \
        to implement the top-priority growth features.\n\n\
        The prompt should:\n\
        1. Be specific to the codebase's tech stack ({lang}{fw})\n\
        2. Reference actual file paths from the manifest\n\
        3. Include step-by-step implementation instructions\n\
        4. Focus on the highest-priority growth opportunities\n\
        5. Be copy-paste ready\n\n\
        ## Manifest\n{manifest_summary}",
        lang = manifest.tech_stack.language,
        fw = manifest.tech_stack.framework.as_deref().map(|f| format!(", {}", f)).unwrap_or_default(),
        manifest_summary = manifest_summary,
    );

    let build_content = llm_client.generate_content(&prompt).await
        .map_err(|e| anyhow::anyhow!(e))?;

    on_progress("build", "Writing implementation prompt...".to_string(), 0.9);

    let output_path = output_dir.join("implementation-prompt.md");
    write_file(&output_path, &build_content).await?;

    on_progress("build", "Implementation prompt generated".to_string(), 1.0);

    let output = EngineOutput::Result {
        manifest_path: Some(manifest_path.to_string_lossy().to_string()),
        template_path: Some(output_path.to_string_lossy().to_string()),
        docs_path: None,
        plan_path: None,
    };
    println!("{}", serde_json::to_string(&output).unwrap());

    Ok(())
}

async fn run_status(input: EngineInput) -> Result<()> {
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

    on_progress("status", "Loading manifest...".to_string(), 0.1);

    let output_dir = PathBuf::from(&input.output_dir);
    let manifest_path = input.manifest_path
        .map(PathBuf::from)
        .unwrap_or_else(|| output_dir.join("growth-manifest.json"));

    let manifest_str = tokio::fs::read_to_string(&manifest_path).await
        .map_err(|_| anyhow::anyhow!("No manifest found at {}. Run analysis first.", manifest_path.display()))?;
    let manifest: GrowthManifest = serde_json::from_str(&manifest_str)?;

    on_progress("status", "Scanning codebase for implementations...".to_string(), 0.3);

    let base_dir = PathBuf::from(&input.project_dir);
    let explorer = CodebaseExplorer::new(base_dir.clone(), Some(input.exclude_folders));

    // Check which growth opportunities have been implemented by scanning the codebase
    let llm_client = create_llm_client(&input.provider, &input.api_key, &input.model, input.base_url.as_deref())
        .map_err(|e| anyhow::anyhow!(e))?;

    on_progress("status", "Checking growth loop implementation status...".to_string(), 0.5);

    let manifest_summary = format_manifest_for_prompt(&manifest);

    // Get file tree for context
    let file_tree = explorer.get_directory_tree(".", 4).await.unwrap_or_default();

    let prompt = format!(
        "You are a growth engineering auditor. Given the following growth manifest and current codebase structure, \
        determine which growth features and opportunities have been implemented, which are in progress, and which are missing.\n\n\
        Format your response as a status report:\n\n\
        ## Growth Loop Status Report\n\n\
        ### Implemented\n- [feature]: [evidence from file tree]\n\n\
        ### In Progress\n- [feature]: [partial evidence]\n\n\
        ### Not Started\n- [feature]: [what's needed]\n\n\
        ### Summary\n- Total features: X\n- Implemented: X\n- In progress: X\n- Not started: X\n\n\
        ## Manifest\n{manifest_summary}\n\n\
        ## Current File Tree\n{file_tree}",
        manifest_summary = manifest_summary,
        file_tree = file_tree,
    );

    let status_content = llm_client.generate_content(&prompt).await
        .map_err(|e| anyhow::anyhow!(e))?;

    on_progress("status", "Writing status report...".to_string(), 0.9);

    let output_path = output_dir.join("growth-status.md");
    write_file(&output_path, &status_content).await?;

    on_progress("status", "Status check complete".to_string(), 1.0);

    let output = EngineOutput::Result {
        manifest_path: Some(manifest_path.to_string_lossy().to_string()),
        template_path: None,
        docs_path: Some(output_path.to_string_lossy().to_string()),
        plan_path: None,
    };
    println!("{}", serde_json::to_string(&output).unwrap());

    Ok(())
}

fn format_manifest_for_prompt(manifest: &GrowthManifest) -> String {
    let mut lines = Vec::new();
    lines.push(format!("**Project:** {}", manifest.project_name));
    if let Some(desc) = &manifest.description {
        lines.push(format!("**Description:** {}", desc));
    }
    lines.push(format!("\n**Tech Stack:** {} {}", 
        manifest.tech_stack.language,
        manifest.tech_stack.framework.as_deref().unwrap_or("")));
    if let Some(db) = &manifest.tech_stack.database {
        lines.push(format!("**Database:** {}", db));
    }
    
    if !manifest.current_growth_features.is_empty() {
        lines.push(format!("\n**Existing Growth Features ({}):**", manifest.current_growth_features.len()));
        for feat in &manifest.current_growth_features {
            lines.push(format!("- {} (in {}, confidence: {:.0}%)", 
                feat.feature_name, feat.file_path, feat.confidence_score * 100.0));
        }
    }
    
    if !manifest.growth_opportunities.is_empty() {
        lines.push(format!("\n**Growth Opportunities ({}):**", manifest.growth_opportunities.len()));
        for opp in &manifest.growth_opportunities {
            lines.push(format!("- [{}] {}: {}", opp.priority.to_uppercase(), opp.feature_name, opp.description));
        }
    }

    if !manifest.revenue_leakage.is_empty() {
        lines.push(format!("\n**Revenue Leakage ({}):**", manifest.revenue_leakage.len()));
        for leak in &manifest.revenue_leakage {
            lines.push(format!("- [{}] {}: {}", leak.impact.to_uppercase(), leak.issue, leak.recommendation));
        }
    }
    
    lines.join("\n")
}
