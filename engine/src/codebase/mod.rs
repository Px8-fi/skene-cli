use std::path::{Path, PathBuf};
use std::collections::HashSet;
use anyhow::{Result, anyhow};
use path_clean::PathClean;
use walkdir::WalkDir;
use glob::glob;
use tokio::fs;

pub mod filters;
pub mod tree;

use self::filters::DEFAULT_EXCLUDE_FOLDERS;

#[derive(Clone)]
pub struct CodebaseExplorer {
    base_dir: PathBuf,
    exclude_folders: HashSet<String>,
}

impl CodebaseExplorer {
    pub fn new(base_dir: PathBuf, exclude_folders: Option<Vec<String>>) -> Self {
        let mut excludes = HashSet::new();
        for folder in DEFAULT_EXCLUDE_FOLDERS {
            excludes.insert(folder.to_string());
        }
        if let Some(custom_excludes) = exclude_folders {
            for folder in custom_excludes {
                excludes.insert(folder);
            }
        }

        Self {
            base_dir: base_dir.clean(),
            exclude_folders: excludes,
        }
    }

    fn resolve_safe_path(&self, relative_path: &str) -> Result<PathBuf> {
        let clean_path = PathBuf::from(relative_path.trim_start_matches('/'));
        let full_path = self.base_dir.join(&clean_path).clean();

        if !full_path.starts_with(&self.base_dir) {
            return Err(anyhow!("Access denied: Path is outside allowed directory"));
        }
        Ok(full_path)
    }

    fn should_exclude(&self, path: &Path) -> bool {
        let path_str = path.to_string_lossy();
        
        // Check if any excluded folder is in the path
        for excluded in &self.exclude_folders {
            // Check exact match in path components
            for component in path.components() {
                if let Some(comp_str) = component.as_os_str().to_str() {
                     if comp_str == excluded {
                         return true;
                     }
                }
            }
            
            // If excluded contains a slash, check as substring
            if (excluded.contains('/') || excluded.contains('\\')) && path_str.contains(excluded) {
                return true;
            }
        }
        false
    }

    pub async fn read_file(&self, file_path: &str) -> Result<String> {
        let target_file = self.resolve_safe_path(file_path)?;
        
        if !target_file.exists() {
            return Err(anyhow!("File does not exist: {}", file_path));
        }
        if !target_file.is_file() {
            return Err(anyhow!("Path is not a file: {}", file_path));
        }

        let content = fs::read_to_string(target_file).await?;
        Ok(content)
    }

    pub async fn list_directory(&self, path: &str) -> Result<Vec<String>> {
        let target_path = self.resolve_safe_path(path)?;
        
        if !target_path.exists() {
            return Err(anyhow!("Path does not exist: {}", path));
        }
        if !target_path.is_dir() {
            return Err(anyhow!("Path is not a directory: {}", path));
        }

        let mut items = Vec::new();
        let mut read_dir = fs::read_dir(target_path).await?;
        
        while let Ok(Some(entry)) = read_dir.next_entry().await {
            let path = entry.path();
            if self.should_exclude(&path) {
                continue;
            }
            
            if let Ok(relative) = path.strip_prefix(&self.base_dir) {
                items.push(relative.to_string_lossy().to_string());
            }
        }
        
        Ok(items)
    }

    pub async fn search_files(&self, pattern: &str) -> Result<Vec<String>> {
        // pattern is a glob pattern relative to base_dir
        // e.g. "**/*.py"
        
        let full_pattern = self.base_dir.join(pattern);
        let pattern_str = full_pattern.to_string_lossy();
        
        let mut matches = Vec::new();
        
        for entry in glob(&pattern_str)? {
            match entry {
                Ok(path) => {
                    if self.should_exclude(&path) {
                        continue;
                    }
                    if path.is_file() {
                        if let Ok(relative) = path.strip_prefix(&self.base_dir) {
                            matches.push(relative.to_string_lossy().to_string());
                        }
                    }
                },
                Err(e) => println!("Glob error: {:?}", e),
            }
        }
        
        Ok(matches)
    }

    pub async fn get_directory_tree(&self, _path: &str, _max_depth: usize) -> Result<String> {
        // Delegate to tree module or implement here
        // For simplicity, implement a basic version using WalkDir
        let target_path = self.resolve_safe_path(_path)?;
        
        let mut output = String::new();
        output.push_str(&format!("{}/\n", target_path.file_name().unwrap_or_default().to_string_lossy()));

        for entry in WalkDir::new(&target_path)
            .min_depth(1)
            .max_depth(_max_depth)
            .sort_by_file_name()
            .into_iter()
            .filter_entry(|e| !self.should_exclude(e.path())) 
        {
            let entry = entry?;
            let depth = entry.depth();
            let indent = "  ".repeat(depth);
            let name = entry.file_name().to_string_lossy();
            if entry.file_type().is_dir() {
                output.push_str(&format!("{}{}/\n", indent, name));
            } else {
                output.push_str(&format!("{}{}\n", indent, name));
            }
        }
        
        Ok(output)
    }
}
