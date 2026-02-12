use crate::strategies::MultiStepStrategy;
use crate::strategies::steps::{select_files::SelectFilesStep, read_files::ReadFilesStep, analyze::AnalyzeStep};
use crate::analyzers::prompts::INDUSTRY_PROMPT;

pub fn create_industry_analyzer() -> MultiStepStrategy {
    MultiStepStrategy::new(vec![
        Box::new(SelectFilesStep::new(
            "Select documentation and metadata files to classify the industry.",
            vec![
                "README.md", "package.json", "pyproject.toml", "Cargo.toml",
                "go.mod", "setup.py", "requirements.txt",
                "docs/**/*", "website/**/*"
            ],
            10,
            "industry_files"
        )),
        Box::new(ReadFilesStep::new("industry_files", "industry_contents")),
        Box::new(AnalyzeStep::new(INDUSTRY_PROMPT, "industry", Some("industry_contents"))),
    ])
}
