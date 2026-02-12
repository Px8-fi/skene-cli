use async_trait::async_trait;
use reqwest::Client;
use serde_json::{json, Value};
use std::error::Error;
use super::LLMClient;
use std::time::Duration;

pub struct OpenAICompatClient {
    client: Client,
    provider: String,
    model: String,
    api_key: String,
    base_url: String,
}

impl OpenAICompatClient {
    pub fn new(provider: &str, model: &str, api_key: &str, base_url: Option<&str>) -> Self {
        let base_url = base_url.map(|s| s.to_string()).unwrap_or_else(|| {
            match provider.to_lowercase().as_str() {
                "openai" => "https://api.openai.com/v1/chat/completions".to_string(),
                "gemini" => "https://generativelanguage.googleapis.com/v1beta/openai/chat/completions".to_string(),
                "ollama" => "http://localhost:11434/v1/chat/completions".to_string(),
                "lmstudio" => "http://localhost:1234/v1/chat/completions".to_string(),
                _ => "https://api.openai.com/v1/chat/completions".to_string(),
            }
        });

        Self {
            client: Client::builder()
                .timeout(Duration::from_secs(120))
                .build()
                .unwrap_or_default(),
            provider: provider.to_string(),
            model: model.to_string(),
            api_key: api_key.to_string(),
            base_url,
        }
    }
}

#[async_trait]
impl LLMClient for OpenAICompatClient {
    async fn generate_content(&self, prompt: &str) -> Result<String, Box<dyn Error + Send + Sync>> {
        let mut request_builder = self.client.post(&self.base_url)
            .header("Content-Type", "application/json");

        // Handle auth based on provider
        if !self.api_key.is_empty() {
             request_builder = request_builder.header("Authorization", format!("Bearer {}", self.api_key));
        }

        let body = json!({
            "model": self.model,
            "messages": [
                {
                    "role": "user",
                    "content": prompt
                }
            ],
            "temperature": 0.7
        });

        let response = request_builder
            .json(&body)
            .send()
            .await?;

        if !response.status().is_success() {
            let error_text = response.text().await?;
            return Err(format!("API request failed: {}", error_text).into());
        }

        let json: Value = response.json().await?;
        
        let content = json["choices"][0]["message"]["content"]
            .as_str()
            .ok_or("No content in response")?
            .to_string();

        Ok(content)
    }

    fn get_model_name(&self) -> String {
        self.model.clone()
    }

    fn get_provider_name(&self) -> String {
        self.provider.clone()
    }
}
