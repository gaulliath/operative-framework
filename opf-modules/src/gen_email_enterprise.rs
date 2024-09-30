use async_trait::async_trait;
use std::collections::HashMap;

use tokio::sync::mpsc::UnboundedSender;

use opf_models::error::{ErrorKind, Module as ModuleError};
use opf_models::event::Event;
use opf_models::{
    metadata::{Arg, Args},
    Target, TargetType,
};

use crate::CompiledModule;

#[derive(Clone)]
pub struct GenerateEmailEnterprise {}

#[async_trait]
impl CompiledModule for GenerateEmailEnterprise {
    fn name(&self) -> String {
        "gen.email_enterprise".to_string()
    }

    fn author(&self) -> String {
        "Tristan Granier".to_string()
    }

    fn resume(&self) -> String {
        "Generate enterprise email".to_string()
    }

    fn args(&self) -> Vec<Arg> {
        vec![
            Arg::new("target_id", true, false, None),
            Arg::new("target", false, false, None),
        ]
    }

    fn target_type(&self) -> TargetType {
        TargetType::Person
    }

    async fn run(
        &self,
        _group_id: i32,
        target: Target,
        _params: Args,
        _tx: Option<UnboundedSender<Event>>,
    ) -> Result<Vec<Target>, ErrorKind> {
        let enterprise =
            target
                .meta
                .get(opf_models::target::ENTERPRISE)
                .ok_or(ErrorKind::Module(ModuleError::Execution(
                    "invalid target, please set a metadata enterprise.".to_string(),
                )))?;

        let first_name =
            target
                .meta
                .get(opf_models::target::FIRST_NAME)
                .ok_or(ErrorKind::Module(ModuleError::Execution(
                    "invalid target, please set a metadata enterprise.".to_string(),
                )))?;

        let last_name = target
            .meta
            .get(opf_models::target::LAST_NAME)
            .ok_or(ErrorKind::Module(ModuleError::Execution(
                "invalid target, please set a metadata enterprise.".to_string(),
            )))?;

        let email = format!(
            "{}.{}@{}.fr",
            first_name.to_lowercase(),
            last_name.to_lowercase(),
            enterprise.to_lowercase()
        );

        let mut results = vec![];

        let mut new_target = HashMap::new();
        new_target.insert(String::from(opf_models::target::NAME), email);
        new_target.insert(
            String::from(opf_models::target::TYPE),
            TargetType::Email.to_string(),
        );
        new_target.insert(
            String::from(opf_models::target::PARENT),
            target.target_id.to_string(),
        );
        results.push(Target::try_from(new_target)?);

        Ok(results)
    }
}

impl GenerateEmailEnterprise {
    pub fn new() -> Box<dyn CompiledModule> {
        Box::new(GenerateEmailEnterprise {})
    }
}
