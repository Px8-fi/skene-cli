use async_trait::async_trait;
use anyhow::Result;
use crate::codebase::CodebaseExplorer;
use crate::llm::LLMClient;
use crate::strategies::context::AnalysisContext;

pub mod select_files;
pub mod read_files;
pub mod analyze;

#[async_trait]
pub trait AnalysisStep: Send + Sync {
    fn name(&self) -> &str;
    async fn execute(&self, codebase: &CodebaseExplorer, llm: &dyn LLMClient, context: &mut AnalysisContext) -> Result<()>;
}
