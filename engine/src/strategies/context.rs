use serde::{Serialize, Deserialize};
use serde_json::Value;
use std::collections::HashMap;

#[derive(Debug, Serialize, Deserialize, Default, Clone)]
pub struct AnalysisMetadata {
    pub model_name: String,
    pub provider_name: String,
    pub total_steps: usize,
    pub files_read: Vec<String>,
    pub tokens_used: usize,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct AnalysisContext {
    pub request: String,
    pub data: HashMap<String, Value>,
    pub metadata: AnalysisMetadata,
}

impl AnalysisContext {
    pub fn new(request: String) -> Self {
        Self {
            request,
            data: HashMap::new(),
            metadata: AnalysisMetadata::default(),
        }
    }

    pub fn set(&mut self, key: &str, value: Value) {
        self.data.insert(key.to_string(), value);
    }

    pub fn get(&self, key: &str) -> Option<&Value> {
        self.data.get(key)
    }
}
