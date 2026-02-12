use anyhow::Result;
use chrono::Local;
use crate::llm::LLMClient;
use crate::manifest::GrowthManifest;
use crate::planner::prompts::{COUNCIL_MEMO_PROMPT_TEMPLATE, ONBOARDING_MEMO_PROMPT_TEMPLATE};

pub mod prompts;

pub struct Planner;

impl Planner {
    pub async fn generate_council_memo(
        llm: &dyn LLMClient,
        manifest: &GrowthManifest,
    ) -> Result<String> {
        let manifest_summary = format_manifest_summary(manifest);
        // Placeholder for template/growth loops until we have structs for them
        let template_section = "";
        let growth_loops_section = "";
        
        let current_time = Local::now().to_rfc3339();
        
        let prompt = COUNCIL_MEMO_PROMPT_TEMPLATE
            .replace("{current_time}", &current_time)
            .replace("{manifest_summary}", &manifest_summary)
            .replace("{template_section}", template_section)
            .replace("{growth_loops_section}", growth_loops_section);
            
        llm.generate_content(&prompt).await.map_err(|e| anyhow::anyhow!(e))
    }

    pub async fn generate_onboarding_memo(
        llm: &dyn LLMClient,
        manifest: &GrowthManifest,
    ) -> Result<String> {
        let manifest_summary = format_manifest_summary(manifest);
        let current_time = Local::now().to_rfc3339();
        
        let prompt = ONBOARDING_MEMO_PROMPT_TEMPLATE
            .replace("{current_time}", &current_time)
            .replace("{manifest_summary}", &manifest_summary)
            .replace("{template_section}", "")
            .replace("{growth_loops_section}", "");
            
        llm.generate_content(&prompt).await.map_err(|e| anyhow::anyhow!(e))
    }
}

fn format_manifest_summary(manifest: &GrowthManifest) -> String {
    let mut lines = Vec::new();
    lines.push(format!("**Project:** {}", manifest.project_name));
    if let Some(desc) = &manifest.description {
        lines.push(format!("**Description:** {}", desc));
    }
    
    lines.push("\n**Tech Stack:**".to_string());
    lines.push(format!("- Language: {}", manifest.tech_stack.language));
    if let Some(fw) = &manifest.tech_stack.framework { lines.push(format!("- Framework: {}", fw)); }
    if let Some(db) = &manifest.tech_stack.database { lines.push(format!("- Database: {}", db)); }
    if let Some(auth) = &manifest.tech_stack.auth { lines.push(format!("- Auth: {}", auth)); }
    
    if !manifest.current_growth_features.is_empty() {
        lines.push(format!("\n**Existing Growth Features:** {} detected", manifest.current_growth_features.len()));
        for feat in manifest.current_growth_features.iter().take(3) {
            lines.push(format!("- {}", feat.feature_name));
        }
    }
    
    if !manifest.growth_opportunities.is_empty() {
        lines.push(format!("\n**Growth Opportunities:** {} identified", manifest.growth_opportunities.len()));
        for opp in manifest.growth_opportunities.iter().take(3) {
            lines.push(format!("- {} (priority: {})", opp.feature_name, opp.priority));
        }
    }
    
    lines.join("\n")
}
