use crate::strategies::MultiStepStrategy;
use crate::strategies::steps::{select_files::SelectFilesStep, read_files::ReadFilesStep, analyze::AnalyzeStep};
use crate::analyzers::prompts::TECH_STACK_PROMPT;

pub fn create_tech_stack_analyzer() -> MultiStepStrategy {
    MultiStepStrategy::new(vec![
        Box::new(SelectFilesStep::new(
            "Select configuration files and representative source files that reveal the technology stack.",
            vec![
                "package.json", "requirements.txt", "pyproject.toml", "Cargo.toml", "go.mod",
                "Gemfile", "composer.json", "*.config.js", "*.config.ts", "tsconfig.json",
                "next.config.*", "vite.config.*", "docker-compose.yml", "Dockerfile",
                ".env.example", "vercel.json", "netlify.toml", "fly.toml", "render.yaml",
                "**/*.py", "**/*.js", "**/*.ts", "**/*.tsx", "**/*.go", "**/*.rs", "**/*.rb"
            ],
            15,
            "tech_stack_files"
        )),
        Box::new(ReadFilesStep::new("tech_stack_files", "tech_stack_contents")),
        Box::new(AnalyzeStep::new(TECH_STACK_PROMPT, "tech_stack", Some("tech_stack_contents"))),
    ])
}
