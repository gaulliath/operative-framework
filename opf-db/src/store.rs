use std::collections::HashMap;

use tokio::sync::mpsc::UnboundedSender;
use tokio::sync::RwLock;

use opf_models::event::Event;
use opf_models::{Link, Target};

#[derive(Debug)]
pub struct DB {
    pub db_tx: UnboundedSender<Event>,
    pub db_id: i32,
    pub targets: RwLock<HashMap<i32, Target>>,
    pub links: RwLock<HashMap<i32, Link>>,
}

impl DB {
    pub fn new(db_tx: UnboundedSender<Event>) -> Self {
        Self {
            db_tx,
            db_id: 0,
            targets: RwLock::new(HashMap::new()),
            links: RwLock::new(HashMap::new()),
        }
    }
}
