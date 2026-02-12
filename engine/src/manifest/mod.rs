use serde::{Deserialize, Serialize};
use chrono::{DateTime, Local};

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct TechStack {
    #[serde(default)]
    pub framework: Option<String>,
    pub language: String,
    #[serde(default)]
    pub database: Option<String>,
    #[serde(default)]
    pub auth: Option<String>,
    #[serde(default)]
    pub deployment: Option<String>,
    #[serde(default)]
    pub package_manager: Option<String>,
    #[serde(default)]
    pub services: Vec<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct GrowthFeature {
    pub feature_name: String,
    pub file_path: String,
    pub detected_intent: String,
    pub confidence_score: f64,
    #[serde(default)]
    pub entry_point: Option<String>,
    #[serde(default)]
    pub growth_potential: Vec<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct GrowthOpportunity {
    pub feature_name: String,
    pub description: String,
    pub priority: String, // "high", "medium", "low"
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct RevenueLeakage {
    pub issue: String,
    #[serde(default)]
    pub file_path: Option<String>,
    pub impact: String, // "high", "medium", "low"
    pub recommendation: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct IndustryInfo {
    #[serde(default)]
    pub primary: Option<String>,
    #[serde(default)]
    pub secondary: Vec<String>,
    #[serde(default)]
    pub confidence: Option<f64>,
    #[serde(default)]
    pub evidence: Vec<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct ProductOverview {
    #[serde(default)]
    pub tagline: Option<String>,
    #[serde(default)]
    pub value_proposition: Option<String>,
    #[serde(default)]
    pub target_audience: Option<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Feature {
    pub name: String,
    pub description: String,
    #[serde(default)]
    pub file_path: Option<String>,
    #[serde(default)]
    pub usage_example: Option<String>,
    #[serde(default)]
    pub category: Option<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct GrowthManifest {
    #[serde(default = "default_version")]
    pub version: String,
    pub project_name: String,
    #[serde(default)]
    pub description: Option<String>,
    pub tech_stack: TechStack,
    #[serde(default)]
    pub industry: Option<IndustryInfo>,
    #[serde(default)]
    pub current_growth_features: Vec<GrowthFeature>,
    #[serde(default)]
    pub growth_opportunities: Vec<GrowthOpportunity>,
    #[serde(default)]
    pub revenue_leakage: Vec<RevenueLeakage>,
    #[serde(default = "default_generated_at")]
    pub generated_at: DateTime<Local>,
}

fn default_version() -> String {
    "1.0".to_string()
}

fn default_generated_at() -> DateTime<Local> {
    Local::now()
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct DocsManifest {
    #[serde(default = "default_docs_version")]
    pub version: String,
    pub project_name: String,
    #[serde(default)]
    pub description: Option<String>,
    pub tech_stack: TechStack,
    #[serde(default)]
    pub product_overview: Option<ProductOverview>,
    #[serde(default)]
    pub industry: Option<IndustryInfo>,
    #[serde(default)]
    pub features: Vec<Feature>,
    #[serde(default)]
    pub current_growth_features: Vec<GrowthFeature>,
    #[serde(default)]
    pub growth_opportunities: Vec<GrowthOpportunity>,
    #[serde(default = "default_generated_at")]
    pub generated_at: DateTime<Local>,
}

fn default_docs_version() -> String {
    "2.0".to_string()
}

// Convert GrowthManifest to DocsManifest
impl From<GrowthManifest> for DocsManifest {
    fn from(manifest: GrowthManifest) -> Self {
        DocsManifest {
            version: "2.0".to_string(),
            project_name: manifest.project_name,
            description: manifest.description,
            tech_stack: manifest.tech_stack,
            product_overview: None,
            industry: manifest.industry,
            features: vec![],
            current_growth_features: manifest.current_growth_features,
            growth_opportunities: manifest.growth_opportunities,
            generated_at: manifest.generated_at,
        }
    }
}
