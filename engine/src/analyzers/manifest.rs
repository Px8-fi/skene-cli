use crate::strategies::MultiStepStrategy;
use crate::strategies::steps::analyze::AnalyzeStep;
use crate::analyzers::prompts::MANIFEST_PROMPT;

pub fn create_manifest_analyzer() -> MultiStepStrategy {
    MultiStepStrategy::new(vec![
        Box::new(AnalyzeStep::new_with_keys(
            MANIFEST_PROMPT, 
            "manifest", 
            vec!["tech_stack", "current_growth_features", "revenue_leakage", "industry"]
        )),
    ])
}
