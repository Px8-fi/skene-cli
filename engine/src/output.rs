use std::path::Path;
use anyhow::Result;
use tokio::fs;
use serde_json::{to_string_pretty, Value};
use crate::manifest::GrowthManifest;

pub async fn write_manifest(path: &Path, manifest: &GrowthManifest) -> Result<()> {
    let json = to_string_pretty(manifest)?;
    if let Some(parent) = path.parent() {
        fs::create_dir_all(parent).await?;
    }
    fs::write(path, json).await?;
    Ok(())
}

pub async fn write_manifest_json(path: &Path, json: &Value) -> Result<()> {
    let content = to_string_pretty(json)?;
    if let Some(parent) = path.parent() {
        fs::create_dir_all(parent).await?;
    }
    fs::write(path, content).await?;
    Ok(())
}

pub async fn write_file(path: &Path, content: &str) -> Result<()> {
    if let Some(parent) = path.parent() {
        fs::create_dir_all(parent).await?;
    }
    fs::write(path, content).await?;
    Ok(())
}

pub async fn write_product_docs(path: &Path, manifest: &Value) -> Result<()> {
    let mut md = String::new();
    
    let name = manifest["project_name"].as_str().unwrap_or("Project");
    md.push_str(&format!("# {}\n\n", name));
    
    if let Some(desc) = manifest["description"].as_str() {
        md.push_str(&format!("{}\n\n", desc));
    }
    
    if let Some(overview) = manifest.get("product_overview") {
        md.push_str("## Product Overview\n\n");
        if let Some(tagline) = overview["tagline"].as_str() {
            md.push_str(&format!("**Tagline:** {}\n\n", tagline));
        }
        if let Some(vp) = overview["value_proposition"].as_str() {
            md.push_str(&format!("**Value Proposition:** {}\n\n", vp));
        }
        if let Some(target) = overview["target_audience"].as_str() {
            md.push_str(&format!("**Target Audience:** {}\n\n", target));
        }
    }
    
    if let Some(features) = manifest["features"].as_array() {
        md.push_str("## Features\n\n");
        for feat in features {
            if let Some(name) = feat["name"].as_str() {
                md.push_str(&format!("### {}\n", name));
                if let Some(desc) = feat["description"].as_str() {
                    md.push_str(&format!("{}\n", desc));
                }
                if let Some(cat) = feat["category"].as_str() {
                    md.push_str(&format!("\n*Category: {}*\n", cat));
                }
                md.push_str("\n");
            }
        }
    }
    
    if let Some(stack) = manifest.get("tech_stack") {
        md.push_str("## Tech Stack\n\n");
        if let Some(lang) = stack["language"].as_str() {
            md.push_str(&format!("- **Language:** {}\n", lang));
        }
        if let Some(fw) = stack["framework"].as_str() {
            md.push_str(&format!("- **Framework:** {}\n", fw));
        }
        // ... more fields ...
    }

    write_file(path, &md).await
}
