use async_trait::async_trait;
use anyhow::Result;
use serde_json::json;
use crate::codebase::CodebaseExplorer;
use crate::llm::LLMClient;
use crate::strategies::context::AnalysisContext;
use crate::strategies::steps::AnalysisStep;

pub struct ReadFilesStep {
    pub source_key: String,
    pub output_key: String,
}

impl ReadFilesStep {
    pub fn new(source_key: &str, output_key: &str) -> Self {
        Self {
            source_key: source_key.to_string(),
            output_key: output_key.to_string(),
        }
    }
}

#[async_trait]
impl AnalysisStep for ReadFilesStep {
    fn name(&self) -> &str {
        "Read Files"
    }

    async fn execute(&self, codebase: &CodebaseExplorer, _llm: &dyn LLMClient, context: &mut AnalysisContext) -> Result<()> {
        let files_value = context.get(&self.source_key).ok_or_else(|| anyhow::anyhow!("Key not found: {}", self.source_key))?;
        let files: Vec<String> = serde_json::from_value(files_value.clone())?;
        
        let mut results = Vec::new();
        
        for file_path in files {
            match codebase.read_file(&file_path).await {
                Ok(content) => {
                    // Truncate if too large? For now, keep it all or simple truncation.
                    let truncated = if content.len() > 50000 {
                        format!("{}... (truncated)", &content[..50000])
                    } else {
                        content
                    };
                    
                    results.push(json!({
                        "path": file_path,
                        "content": truncated
                    }));
                },
                Err(e) => {
                    results.push(json!({
                        "path": file_path,
                        "error": e.to_string()
                    }));
                }
            }
        }

        context.set(&self.output_key, json!(results));
        Ok(())
    }
}
