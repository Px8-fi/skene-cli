use serde::{Deserialize, Deserializer, Serialize};
use serde_json::Value;
use chrono::{DateTime, Local, NaiveDateTime};

/// Deserialize a Vec that may be null in JSON as an empty Vec
fn null_as_empty_vec<'de, D, T>(deserializer: D) -> Result<Vec<T>, D::Error>
where
    D: Deserializer<'de>,
    T: Deserialize<'de>,
{
    let opt: Option<Vec<T>> = Option::deserialize(deserializer)?;
    Ok(opt.unwrap_or_default())
}

/// Deserialize a Vec<T> that may be either a JSON array OR a JSON object (map).
/// When the LLM returns a map like {"key1": {...}, "key2": {...}} instead of
/// an array [{...}, {...}], we extract the map values and deserialize each as T.
/// Also handles null as an empty Vec.
fn flexible_vec<'de, D, T>(deserializer: D) -> Result<Vec<T>, D::Error>
where
    D: Deserializer<'de>,
    T: Deserialize<'de>,
{
    let value: Option<Value> = Option::deserialize(deserializer)?;
    match value {
        None => Ok(Vec::new()),
        Some(Value::Array(arr)) => {
            let mut result = Vec::new();
            for item in arr {
                match serde_json::from_value(item) {
                    Ok(v) => result.push(v),
                    Err(_) => {} // skip items that don't match the schema
                }
            }
            Ok(result)
        }
        Some(Value::Object(map)) => {
            let mut result = Vec::new();
            for (_key, item) in map {
                match serde_json::from_value(item) {
                    Ok(v) => result.push(v),
                    Err(_) => {} // skip items that don't match the schema
                }
            }
            Ok(result)
        }
        Some(_) => Ok(Vec::new()),
    }
}

/// Deserialize a datetime that may be naive (no timezone), RFC 3339, null, or missing
fn flexible_datetime<'de, D>(deserializer: D) -> Result<DateTime<Local>, D::Error>
where
    D: Deserializer<'de>,
{
    let opt: Option<String> = Option::deserialize(deserializer)?;
    let s = match opt {
        Some(s) if !s.is_empty() => s,
        _ => return Ok(Local::now()),
    };
    // Try RFC 3339 first (with timezone)
    if let Ok(dt) = DateTime::parse_from_rfc3339(&s) {
        return Ok(dt.with_timezone(&Local));
    }
    // Try common formats without timezone
    for fmt in &["%Y-%m-%dT%H:%M:%S%.f", "%Y-%m-%dT%H:%M:%S", "%Y-%m-%d %H:%M:%S"] {
        if let Ok(naive) = NaiveDateTime::parse_from_str(&s, fmt) {
            return Ok(naive.and_local_timezone(Local).unwrap());
        }
    }
    // Fall back to now instead of failing
    Ok(Local::now())
}

#[derive(Debug, Serialize, Deserialize, Clone, Default)]
pub struct TechStack {
    #[serde(default)]
    pub framework: Option<String>,
    #[serde(default)]
    pub language: String,
    #[serde(default)]
    pub database: Option<String>,
    #[serde(default)]
    pub auth: Option<String>,
    #[serde(default)]
    pub deployment: Option<String>,
    #[serde(default)]
    pub package_manager: Option<String>,
    #[serde(default, deserialize_with = "null_as_empty_vec")]
    pub services: Vec<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct GrowthFeature {
    #[serde(default)]
    pub feature_name: String,
    #[serde(default)]
    pub file_path: String,
    #[serde(default)]
    pub detected_intent: String,
    #[serde(default)]
    pub confidence_score: f64,
    #[serde(default)]
    pub entry_point: Option<String>,
    #[serde(default, deserialize_with = "null_as_empty_vec")]
    pub growth_potential: Vec<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct GrowthOpportunity {
    #[serde(default, alias = "opportunity")]
    pub feature_name: String,
    #[serde(default)]
    pub description: String,
    #[serde(default)]
    pub priority: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct RevenueLeakage {
    #[serde(default)]
    pub issue: String,
    #[serde(default)]
    pub file_path: Option<String>,
    #[serde(default)]
    pub impact: String,
    #[serde(default)]
    pub recommendation: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct IndustryInfo {
    #[serde(default)]
    pub primary: Option<String>,
    #[serde(default, deserialize_with = "null_as_empty_vec")]
    pub secondary: Vec<String>,
    #[serde(default)]
    pub confidence: Option<f64>,
    #[serde(default, deserialize_with = "null_as_empty_vec")]
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
    #[serde(default)]
    pub name: String,
    #[serde(default)]
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
    #[serde(default)]
    pub project_name: String,
    #[serde(default)]
    pub description: Option<String>,
    #[serde(default)]
    pub tech_stack: TechStack,
    #[serde(default)]
    pub industry: Option<IndustryInfo>,
    #[serde(default, deserialize_with = "flexible_vec")]
    pub current_growth_features: Vec<GrowthFeature>,
    #[serde(default, deserialize_with = "flexible_vec")]
    pub growth_opportunities: Vec<GrowthOpportunity>,
    #[serde(default, deserialize_with = "flexible_vec")]
    pub revenue_leakage: Vec<RevenueLeakage>,
    #[serde(default = "default_generated_at", deserialize_with = "flexible_datetime")]
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
    #[serde(default)]
    pub project_name: String,
    #[serde(default)]
    pub description: Option<String>,
    #[serde(default)]
    pub tech_stack: TechStack,
    #[serde(default)]
    pub product_overview: Option<ProductOverview>,
    #[serde(default)]
    pub industry: Option<IndustryInfo>,
    #[serde(default, deserialize_with = "flexible_vec")]
    pub features: Vec<Feature>,
    #[serde(default, deserialize_with = "flexible_vec")]
    pub current_growth_features: Vec<GrowthFeature>,
    #[serde(default, deserialize_with = "flexible_vec")]
    pub growth_opportunities: Vec<GrowthOpportunity>,
    #[serde(default = "default_generated_at", deserialize_with = "flexible_datetime")]
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
