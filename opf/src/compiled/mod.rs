pub mod args;
pub mod linkedin_search;
pub mod sample;

use crate::common::Target;
use crate::core::session::Session;
use crate::error::ErrorKind;
use crate::modules::metadata::{Arg, Args};
use crate::modules::{Manager, OpfModule};

pub trait Compiled {
    fn name(&self) -> String;
    fn author(&self) -> String;
    fn resume(&self) -> String;
    fn args(&self) -> Vec<Arg>;
    fn run(&self, sess: &mut Session, params: Args) -> Result<Vec<Target>, ErrorKind>;
}

pub fn registers(manager: &mut Manager, modules: Vec<Box<dyn Compiled>>) -> Result<(), ErrorKind> {
    for compiled in modules {
        manager
            .modules
            .insert(compiled.name(), OpfModule::Compiled(compiled));
    }
    Ok(())
}
