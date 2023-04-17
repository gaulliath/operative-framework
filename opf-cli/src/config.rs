use std::sync::{Arc, RwLock};

use lazy_static::lazy_static;

lazy_static! {
    pub static ref WORKSPACE: Arc<RwLock<String>> = Arc::new(RwLock::new("".to_string()));
}
