pub mod context;
pub mod steps;

use anyhow::Result;
use crate::codebase::CodebaseExplorer;
use crate::llm::LLMClient;
use self::context::AnalysisContext;
use self::steps::AnalysisStep;

pub struct MultiStepStrategy {
    steps: Vec<Box<dyn AnalysisStep>>,
}

impl MultiStepStrategy {
    pub fn new(steps: Vec<Box<dyn AnalysisStep>>) -> Self {
        Self { steps }
    }

    pub async fn run<F>(
        &self,
        codebase: &CodebaseExplorer,
        llm: &dyn LLMClient,
        request: String,
        on_progress: F,
    ) -> Result<AnalysisContext>
    where
        F: Fn(String, f64, usize, usize) + Send + Sync,
    {
        self.run_with_context(codebase, llm, request, None, on_progress).await
    }

    pub async fn run_with_context<F>(
        &self,
        codebase: &CodebaseExplorer,
        llm: &dyn LLMClient,
        request: String,
        initial_context: Option<AnalysisContext>,
        on_progress: F,
    ) -> Result<AnalysisContext>
    where
        F: Fn(String, f64, usize, usize) + Send + Sync,
    {
        let mut context = initial_context.unwrap_or_else(|| AnalysisContext::new(request));
        context.metadata.model_name = llm.get_model_name();
        context.metadata.provider_name = llm.get_provider_name();
        context.metadata.total_steps = self.steps.len();

        for (i, step) in self.steps.iter().enumerate() {
            let step_num = i + 1;
            let total = self.steps.len();
            let progress = (i as f64 / total as f64) * 100.0;
            
            on_progress(format!("Executing {}", step.name()), progress, step_num, total);
            
            step.execute(codebase, llm, &mut context).await?;
        }
        
        on_progress("Complete".to_string(), 100.0, self.steps.len(), self.steps.len());

        Ok(context)
    }
}
