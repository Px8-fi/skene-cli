use crate::strategies::MultiStepStrategy;
use crate::strategies::steps::{select_files::SelectFilesStep, read_files::ReadFilesStep, analyze::AnalyzeStep};
use crate::analyzers::prompts::REVENUE_LEAKAGE_PROMPT;

pub fn create_revenue_leakage_analyzer() -> MultiStepStrategy {
    MultiStepStrategy::new(vec![
        Box::new(SelectFilesStep::new(
            "Select source files that might reveal revenue leakage or pricing logic.",
            vec![
                "**/*pricing*", "**/*billing*", "**/*subscription*", "**/*payment*",
                "**/*plan*", "**/*tier*", "**/*upgrade*", "**/*limit*",
                "**/*user*", "**/*auth*", "**/*api/**/*",
            ],
            20,
            "revenue_files"
        )),
        Box::new(ReadFilesStep::new("revenue_files", "revenue_contents")),
        Box::new(AnalyzeStep::new(REVENUE_LEAKAGE_PROMPT, "revenue_leakage", Some("revenue_contents"))),
    ])
}
