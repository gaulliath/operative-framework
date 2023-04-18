use std::collections::HashMap;
use std::net::{SocketAddr, ToSocketAddrs};

use std::time::Duration;

use async_trait::async_trait;

use tokio::net::TcpStream;
use tokio::sync::mpsc::UnboundedSender;

use opf_models::error::{ErrorKind, Module as ModuleError};
use opf_models::event::{send_event, Event};
use opf_models::{
    metadata::{Arg, Args},
    Target, TargetType,
};

use crate::CompiledModule;

const MOST_COMMON_PORTS: &[u16] = &[
    80, 23, 443, 21, 22, 25, 3389, 110, 445, 139, 143, 53, 135, 3306, 8080, 1723, 111, 995, 993,
    5900, 1025, 587, 8888, 199, 1720, 465, 548, 113, 81, 6001, 10000, 514, 5060, 179, 1026, 2000,
    8443, 8000, 32768, 554, 26, 1433, 49152, 2001, 515, 8008, 49154, 1027, 5666, 646, 5000, 5631,
    631, 49153, 8081, 2049, 88, 79, 5800, 106, 2121, 1110, 49155, 6000, 513, 990, 5357, 427, 49156,
    543, 544, 5101, 144, 7, 389, 8009, 3128, 444, 9999, 5009, 7070, 5190, 3000, 5432, 1900, 3986,
    13, 1029, 9, 5051, 6646, 49157, 1028, 873, 1755, 2717, 4899, 9100, 119, 37, 1000, 3001, 5001,
    82, 10010, 1030, 9090, 2107, 1024, 2103, 6004, 1801, 5050, 19, 8031, 1041, 255, 2967, 1049,
    1048, 1053, 3703, 1056, 1065, 1064, 1054, 17, 808, 3689, 1031, 1044, 1071, 5901, 9102, 100,
    8010, 2869, 1039, 5120, 4001, 9000, 2105, 636, 1038, 2601, 7000, 1, 1066, 1069, 625, 311, 280,
    254, 4000, 5003, 1761, 2002, 2005, 1998, 1032, 1050, 6112, 3690, 1521, 2161, 6002, 1080, 2401,
    4045, 902, 7937, 787, 1058, 2383, 32771, 1033, 1040, 1059, 50000, 5555, 10001, 1494, 2301, 593,
    3, 3268, 7938, 1234, 1022, 1035, 9001, 1074, 8002, 1036, 1037, 464, 1935, 6666, 2003, 497,
    5601, 9200, 9300,
];

#[derive(Clone)]
pub struct PortScanner {}

#[async_trait]
impl CompiledModule for PortScanner {
    fn name(&self) -> String {
        "scan.ports".to_string()
    }

    fn author(&self) -> String {
        "Tristan Granier".to_string()
    }

    fn is_threaded(&self) -> bool {
        true
    }

    fn resume(&self) -> String {
        "Checking if common port is open.".to_string()
    }

    fn args(&self) -> Vec<Arg> {
        vec![
            Arg::new("target_id", true, false, None),
            Arg::new("target", false, false, None),
        ]
    }

    async fn run(
        &self,
        params: Args,
        tx: Option<UnboundedSender<Event>>,
    ) -> Result<Vec<Target>, ErrorKind> {
        let company = params.get("target").unwrap();
        let target_id = params
            .get("target_id")
            .ok_or(ErrorKind::Module(ModuleError::TargetNotAvailable))?
            .value
            .ok_or(ErrorKind::Module(ModuleError::TargetNotAvailable))?;
        let target = company
            .value
            .ok_or(ErrorKind::Module(ModuleError::TargetNotAvailable))?;
        let tx = tx.ok_or(ErrorKind::Module(ModuleError::ParamNotAvailable(
            "tx is mandatory for threaded module".to_string(),
        )))?;

        for port in MOST_COMMON_PORTS {
            let port = port.clone();
            let _ = tokio::spawn(PortScanner::scan_port(
                tx.clone(),
                target.clone(),
                port,
                target_id.clone(),
            ));
        }

        Ok(vec![])
    }

    fn target_type(&self) -> TargetType {
        TargetType::Domain
    }
}

impl PortScanner {
    pub fn new() -> Box<dyn CompiledModule> {
        Box::new(PortScanner {})
    }

    async fn scan_port(tx: UnboundedSender<Event>, hostname: String, port: u16, parent: String) {
        let timeout = Duration::from_secs(3);
        let socket_addresses: Vec<SocketAddr> = format!("{}:{}", hostname, port)
            .to_socket_addrs()
            .expect("port scanner: Creating socket address")
            .collect();

        if socket_addresses.len() == 0 {
            return;
        }

        // TODO: detect protocol
        match tokio::time::timeout(timeout, TcpStream::connect(&socket_addresses[0])).await {
            Ok(Ok(_)) => {
                let mut target = HashMap::new();
                target.insert(String::from("name"), port.to_string());
                target.insert(String::from("type"), String::from("port"));
                target.insert(String::from("parent"), parent);

                if let Ok(target) = Target::try_from(target) {
                    let _ = send_event(&tx, Event::ResultsModule(vec![target])).await;
                    let _ = send_event(
                        &tx,
                        Event::ResponseSimple(format!(
                            "hostname '{}' have a new port '{}'",
                            hostname, port
                        )),
                    )
                    .await;
                }
            }
            Ok(Err(e)) => {
                let _ = send_event(&tx, Event::ResponseError(e.to_string()));
            }
            _ => {}
        }
    }
}
