use crate::strategies::MultiStepStrategy;
use crate::strategies::steps::{select_files::SelectFilesStep, read_files::ReadFilesStep, analyze::AnalyzeStep};
use crate::analyzers::prompts::GROWTH_FEATURES_PROMPT;

pub fn create_growth_features_analyzer() -> MultiStepStrategy {
    MultiStepStrategy::new(vec![
        Box::new(SelectFilesStep::new(
            "Select source files that might contain growth-related features.",
            vec![
                "**/*.py", "**/*.ts", "**/*.tsx", "**/*.js", "**/*.jsx",
                "**/routes/**/*", "**/api/**/*", "**/features/**/*",
                "**/components/**/*", "**/pages/**/*", "**/app/**/*",
            ],
            30,
            "growth_files"
        )),
        Box::new(ReadFilesStep::new("growth_files", "growth_contents")),
        Box::new(AnalyzeStep::new(GROWTH_FEATURES_PROMPT, "current_growth_features", Some("growth_contents"))),
    ])
}
