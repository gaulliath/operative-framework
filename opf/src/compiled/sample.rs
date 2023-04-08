use crate::common::Target;
use crate::compiled::Compiled;
use crate::core::session::Session;
use crate::error::ErrorKind;
use crate::modules::metadata::{Arg, Args};

pub struct SampleCompiled {}

impl Compiled for SampleCompiled {
    fn name(&self) -> String {
        "sample.compiled".to_string()
    }

    fn author(&self) -> String {
        "Tristan Granier".to_string()
    }

    fn resume(&self) -> String {
        "Sample module, nothing happened here".to_string()
    }

    fn args(&self) -> Vec<Arg> {
        vec![]
    }

    fn run(&self, sess: &mut Session, params: Args) -> Result<Vec<Target>, ErrorKind> {
        println!("Execution...");
        println!("Session {}", sess.session_id);
        Ok(vec![])
    }
}

impl SampleCompiled {
    pub fn new() -> Box<dyn Compiled> {
        Box::new(SampleCompiled {})
    }
}
