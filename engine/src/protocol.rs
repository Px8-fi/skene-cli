use serde::{Deserialize, Deserializer, Serialize};

/// Deserialize a Vec that may be null in JSON as an empty Vec
fn null_as_empty_vec<'de, D, T>(deserializer: D) -> Result<Vec<T>, D::Error>
where
    D: Deserializer<'de>,
    T: Deserialize<'de>,
{
    let opt: Option<Vec<T>> = Option::deserialize(deserializer)?;
    Ok(opt.unwrap_or_default())
}

#[derive(Debug, Deserialize)]
pub struct EngineInput {
    pub command: String, // analyze, plan, build
    pub provider: String,
    pub model: String,
    pub api_key: String,
    pub base_url: Option<String>,
    pub project_dir: String,
    pub output_dir: String,
    #[serde(default)]
    pub product_docs: bool,
    #[serde(default, deserialize_with = "null_as_empty_vec")]
    pub exclude_folders: Vec<String>,
    #[serde(default)]
    pub debug: bool,
    // For "plan" command
    pub manifest_path: Option<String>,
    pub template_path: Option<String>,
    pub onboarding: Option<bool>,
}

#[derive(Debug, Serialize)]
#[serde(tag = "type", rename_all = "snake_case")]
pub enum EngineOutput {
    Progress {
        phase: String,
        step: usize,
        total_steps: usize,
        progress: f64,
        message: String,
    },
    Result {
        manifest_path: Option<String>,
        template_path: Option<String>,
        docs_path: Option<String>,
        plan_path: Option<String>,
    },
    Error {
        message: String,
        code: Option<String>,
    },
}
