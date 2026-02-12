use async_trait::async_trait;
use std::error::Error;

#[async_trait]
pub trait LLMClient: Send + Sync {
    /// Generate text content from a prompt
    async fn generate_content(&self, prompt: &str) -> Result<String, Box<dyn Error + Send + Sync>>;
    
    /// Get the name of the model being used
    fn get_model_name(&self) -> String;
    
    /// Get the provider name
    fn get_provider_name(&self) -> String;
}

pub mod openai_compat;
pub mod anthropic;
pub mod debug;

use self::openai_compat::OpenAICompatClient;
use self::anthropic::AnthropicClient;

pub fn create_llm_client(
    provider: &str,
    api_key: &str,
    model: &str,
    base_url: Option<&str>,
) -> Result<Box<dyn LLMClient>, String> {
    match provider.to_lowercase().as_str() {
        "openai" | "gemini" | "ollama" | "lmstudio" | "generic" | "openai-compatible" => {
            Ok(Box::new(OpenAICompatClient::new(
                provider,
                model,
                api_key,
                base_url,
            )))
        }
        "anthropic" | "claude" => {
            Ok(Box::new(AnthropicClient::new(
                model,
                api_key,
                base_url,
            )))
        }
        _ => Err(format!("Unknown provider: {}", provider)),
    }
}
