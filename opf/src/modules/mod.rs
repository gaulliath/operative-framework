pub mod metadata;

use crate::error;

#[derive(Debug, Default)]
pub struct Manager {
    pub modules: Vec<Module>,
}

#[derive(Debug, Clone)]
pub struct Module {
    pub file_name: String,
    pub metadata: metadata::Metadata,
}

pub fn register_modules<'a>(
    manager: &'a mut Manager,
    directory: &'a str,
) -> Result<(), error::Manager> {
    let paths = match std::fs::read_dir(directory) {
        Ok(paths) => paths,
        Err(_) => return Err(error::Manager::CantOpenDirectory(directory.to_string())),
    };

    for p in paths {
        let path = match p {
            Ok(path) => path,
            Err(_) => return Err(error::Manager::CantProcessFile),
        };

        let file_name = path.path().display().to_string();

        let content = match std::fs::read_to_string(String::from(&file_name)) {
            Ok(content) => content,
            Err(_) => {
                let name = file_name.clone();
                return Err(error::Manager::CantReadContent(name));
            }
        };

        let metadata = match metadata::parse(content.as_str()) {
            Ok(meta) => meta,
            Err(_) => return Err(error::Manager::CantParseMetadata(String::from(&file_name))),
        };

        manager.modules.push(Module {
            file_name,
            metadata,
        })
    }
    Ok(())
}
