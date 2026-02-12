use crate::strategies::MultiStepStrategy;
use crate::strategies::steps::{select_files::SelectFilesStep, read_files::ReadFilesStep, analyze::AnalyzeStep};
use crate::analyzers::prompts::{PRODUCT_OVERVIEW_PROMPT, FEATURES_PROMPT, DOCS_MANIFEST_PROMPT};

pub fn create_product_overview_analyzer() -> MultiStepStrategy {
    MultiStepStrategy::new(vec![
        Box::new(SelectFilesStep::new(
            "Select documentation files for product overview.",
            vec!["README.md", "package.json", "docs/**/*", "about/**/*"],
            10,
            "overview_files"
        )),
        Box::new(ReadFilesStep::new("overview_files", "overview_contents")),
        Box::new(AnalyzeStep::new(PRODUCT_OVERVIEW_PROMPT, "product_overview", Some("overview_contents"))),
    ])
}

pub fn create_features_analyzer() -> MultiStepStrategy {
    MultiStepStrategy::new(vec![
        Box::new(SelectFilesStep::new(
            "Select source files to document user-facing features.",
            vec![
                "**/*.py", "**/*.ts", "**/*.tsx", "**/*.js", 
                "**/routes/**/*", "**/api/**/*", "**/features/**/*", 
                "**/components/**/*"
            ],
            20,
            "features_files"
        )),
        Box::new(ReadFilesStep::new("features_files", "features_contents")),
        Box::new(AnalyzeStep::new(FEATURES_PROMPT, "features", Some("features_contents"))),
    ])
}

pub fn create_docs_manifest_analyzer() -> MultiStepStrategy {
    MultiStepStrategy::new(vec![
        Box::new(AnalyzeStep::new_with_keys(
            DOCS_MANIFEST_PROMPT, 
            "docs_manifest", 
            vec!["tech_stack", "product_overview", "industry", "features", "current_growth_features"]
        )),
    ])
}
