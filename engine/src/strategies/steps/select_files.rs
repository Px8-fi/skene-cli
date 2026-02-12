use async_trait::async_trait;
use anyhow::Result;
use serde_json::json;
use std::collections::HashSet;
use crate::codebase::CodebaseExplorer;
use crate::llm::LLMClient;
use crate::strategies::context::AnalysisContext;
use crate::strategies::steps::AnalysisStep;

pub struct SelectFilesStep {
    pub prompt: String,
    pub patterns: Vec<String>,
    pub max_files: usize,
    pub output_key: String,
}

impl SelectFilesStep {
    pub fn new(prompt: &str, patterns: Vec<&str>, max_files: usize, output_key: &str) -> Self {
        Self {
            prompt: prompt.to_string(),
            patterns: patterns.iter().map(|s| s.to_string()).collect(),
            max_files,
            output_key: output_key.to_string(),
        }
    }
}

#[async_trait]
impl AnalysisStep for SelectFilesStep {
    fn name(&self) -> &str {
        "Select Files"
    }

    async fn execute(&self, codebase: &CodebaseExplorer, llm: &dyn LLMClient, context: &mut AnalysisContext) -> Result<()> {
        let mut candidates = HashSet::new();
        
        for pattern in &self.patterns {
            let matches = codebase.search_files(pattern).await?;
            for path in matches {
                candidates.insert(path);
            }
        }
        
        let mut candidate_list: Vec<String> = candidates.into_iter().collect();
        candidate_list.sort();

        let selected_files = if candidate_list.len() <= self.max_files {
            candidate_list
        } else {
            // Use LLM to select best files
            let candidate_text = candidate_list.join("\n");
            let prompt = format!(
                "{}\n\nAvailable files:\n{}\n\nSelect the most relevant files (max {}). Return strictly a JSON array of file paths.",
                self.prompt, candidate_text, self.max_files
            );
            
            let response = llm.generate_content(&prompt).await.map_err(|e| anyhow::anyhow!(e))?;
            // Clean response to get JSON
            let json_str = response.trim();
            let json_start = json_str.find('[').unwrap_or(0);
            let json_end = json_str.rfind(']').map(|i| i + 1).unwrap_or(json_str.len());
            let json_part = &json_str[json_start..json_end];
            
            match serde_json::from_str::<Vec<String>>(json_part) {
                Ok(files) => files,
                Err(_) => candidate_list.into_iter().take(self.max_files).collect(),
            }
        };

        context.set(&self.output_key, json!(selected_files));
        Ok(())
    }
}
