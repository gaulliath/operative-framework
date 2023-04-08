pub mod metadata;

use crate::compiled::Compiled;
use crate::error::{ErrorKind, Manager as ManagerError, Metadata as MetadataError};
use crate::modules::metadata::Arg;
use std::collections::HashMap;
use std::fmt::{Debug, Formatter};

#[derive(Default)]
pub struct Manager {
    pub modules: HashMap<String, OpfModule>,
    //pub compiled: HashMap<String, Box<dyn Compiled>>,
}

pub enum OpfModule {
    Lua(Module),
    Compiled(Box<dyn Compiled>),
}

#[derive(Debug, Clone)]
pub struct Module {
    pub file_name: String,
    pub metadata: metadata::Metadata,
}

impl OpfModule {
    pub fn name(&self) -> String {
        match &self {
            Self::Lua(ref module) => module.metadata.name.clone().unwrap_or("-".to_string()),
            Self::Compiled(ref module) => module.name(),
        }
    }

    pub fn author(&self) -> String {
        match &self {
            Self::Lua(ref module) => module.metadata.author.clone().unwrap_or("-".to_string()),
            Self::Compiled(ref module) => module.author(),
        }
    }

    pub fn resume(&self) -> String {
        match &self {
            Self::Lua(ref module) => module
                .metadata
                .description
                .clone()
                .unwrap_or("-".to_string()),
            Self::Compiled(ref module) => module.resume(),
        }
    }

    pub fn args(&self) -> Vec<Arg> {
        match &self {
            Self::Lua(ref module) => module.metadata.args.clone(),
            Self::Compiled(ref module) => module.args(),
        }
    }
}

pub fn register(manager: &mut Manager, directory: &str) -> Result<(), ErrorKind> {
    let paths = std::fs::read_dir(directory)
        .map_err(|_| ErrorKind::Manager(ManagerError::CantOpenDirectory(directory.to_string())))?;

    for p in paths {
        let path = p.map_err(|_| ErrorKind::Manager(ManagerError::CantProcessFile))?;

        let file_name = path.path().display().to_string();

        let content = std::fs::read_to_string(String::from(&file_name))
            .map_err(|_| ErrorKind::Manager(ManagerError::CantReadContent(file_name.clone())))?;

        let metadata = metadata::parse(content.as_str())?;
        if metadata.name.is_none() {
            return Err(ErrorKind::Metadata(MetadataError::Required(
                "name".to_string(),
            )));
        }

        manager.modules.insert(
            metadata.name.clone().unwrap(),
            OpfModule::Lua(Module {
                file_name,
                metadata,
            }),
        );
    }
    Ok(())
}

impl Debug for Manager {
    fn fmt(&self, f: &mut Formatter<'_>) -> std::fmt::Result {
        write!(f, "manager: modules {}", self.modules.len(),)
    }
}
