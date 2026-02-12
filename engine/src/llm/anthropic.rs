use async_trait::async_trait;
use reqwest::Client;
use serde_json::{json, Value};
use std::error::Error;
use super::LLMClient;
use std::time::Duration;

pub struct AnthropicClient {
    client: Client,
    model: String,
    api_key: String,
    base_url: String,
}

impl AnthropicClient {
    pub fn new(model: &str, api_key: &str, base_url: Option<&str>) -> Self {
        let base_url = base_url.map(|s| s.to_string())
            .unwrap_or_else(|| "https://api.anthropic.com/v1/messages".to_string());

        Self {
            client: Client::builder()
                .timeout(Duration::from_secs(120))
                .build()
                .unwrap_or_default(),
            model: model.to_string(),
            api_key: api_key.to_string(),
            base_url,
        }
    }
}

#[async_trait]
impl LLMClient for AnthropicClient {
    async fn generate_content(&self, prompt: &str) -> Result<String, Box<dyn Error + Send + Sync>> {
        let body = json!({
            "model": self.model,
            "max_tokens": 4096,
            "messages": [
                {
                    "role": "user",
                    "content": prompt
                }
            ]
        });

        let response = self.client.post(&self.base_url)
            .header("x-api-key", &self.api_key)
            .header("anthropic-version", "2023-06-01")
            .header("Content-Type", "application/json")
            .json(&body)
            .send()
            .await?;

        if !response.status().is_success() {
            let error_text = response.text().await?;
            return Err(format!("API request failed: {}", error_text).into());
        }

        let json: Value = response.json().await?;
        
        // Anthropic returns content as a list of blocks
        let content = json["content"][0]["text"]
            .as_str()
            .ok_or("No content in response")?
            .to_string();

        Ok(content)
    }

    fn get_model_name(&self) -> String {
        self.model.clone()
    }

    fn get_provider_name(&self) -> String {
        "anthropic".to_string()
    }
}
