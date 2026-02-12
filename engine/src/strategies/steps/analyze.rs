use async_trait::async_trait;
use anyhow::Result;
use serde_json::Value;
use crate::codebase::CodebaseExplorer;
use crate::llm::LLMClient;
use crate::strategies::context::AnalysisContext;
use crate::strategies::steps::AnalysisStep;

pub struct AnalyzeStep {
    pub prompt: String,
    pub output_key: String,
    pub source_keys: Vec<String>,
}

impl AnalyzeStep {
    pub fn new(prompt: &str, output_key: &str, source_key: Option<&str>) -> Self {
        Self {
            prompt: prompt.to_string(),
            output_key: output_key.to_string(),
            source_keys: source_key.map(|s| vec![s.to_string()]).unwrap_or_default(),
        }
    }

    pub fn new_with_keys(prompt: &str, output_key: &str, source_keys: Vec<&str>) -> Self {
        Self {
            prompt: prompt.to_string(),
            output_key: output_key.to_string(),
            source_keys: source_keys.iter().map(|s| s.to_string()).collect(),
        }
    }
}

#[async_trait]
impl AnalysisStep for AnalyzeStep {
    fn name(&self) -> &str {
        "Analyze"
    }

    async fn execute(&self, _codebase: &CodebaseExplorer, llm: &dyn LLMClient, context: &mut AnalysisContext) -> Result<()> {
        let mut final_prompt = self.prompt.clone();
        
        for key in &self.source_keys {
            if let Some(data) = context.get(key) {
                // specific formatting for file_contents
                if key.contains("contents") || key.contains("files") {
                    if let Some(array) = data.as_array() {
                        let mut file_context = String::new();
                        for item in array {
                            if let (Some(path), Some(content)) = (item["path"].as_str(), item["content"].as_str()) {
                                file_context.push_str(&format!("\n--- {} ---\n{}\n", path, content));
                            }
                        }
                        final_prompt.push_str(&format!("\n\nFiles ({}) :\n{}", key, file_context));
                    } else {
                        final_prompt.push_str(&format!("\n\n{} (Context):\n{}", key, serde_json::to_string_pretty(data)?));
                    }
                } else {
                    final_prompt.push_str(&format!("\n\n{} (Analysis Result):\n{}", key, serde_json::to_string_pretty(data)?));
                }
            }
        }

        let response = llm.generate_content(&final_prompt).await.map_err(|e| anyhow::anyhow!(e))?;
        
        // Try to parse as JSON, otherwise string
        let result: Value = match extract_json(&response) {
            Ok(json) => json,
            Err(_) => Value::String(response),
        };

        context.set(&self.output_key, result);
        Ok(())
    }
}

fn extract_json(text: &str) -> Result<Value> {
    let text = text.trim();
    if text.starts_with('{') || text.starts_with('[') {
        return Ok(serde_json::from_str(text)?);
    }
    
    // Try to find Markdown JSON block
    if let Some(start) = text.find("```json") {
        let after_start = &text[start+7..];
        if let Some(end) = after_start.find("```") {
             return Ok(serde_json::from_str(&after_start[..end])?);
        }
    }
    // Try to find Markdown block without json tag
    if let Some(start) = text.find("```") {
        let after_start = &text[start+3..];
        if let Some(end) = after_start.find("```") {
             // check if it looks like json
             let content = &after_start[..end].trim();
             if content.starts_with('{') || content.starts_with('[') {
                 if let Ok(val) = serde_json::from_str(content) {
                     return Ok(val);
                 }
             }
        }
    }
    
    // Try to find first { or [
    if let Some(start) = text.find(|c| c == '{' || c == '[') {
        if let Some(end) = text.rfind(|c| c == '}' || c == ']') {
            if end > start {
                return Ok(serde_json::from_str(&text[start..=end])?);
            }
        }
    }

    Err(anyhow::anyhow!("No JSON found"))
}
