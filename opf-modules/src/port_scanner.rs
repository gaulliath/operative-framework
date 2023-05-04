use crate::CompiledModule;
use async_trait::async_trait;
use opf_models::error::{ErrorKind, Module as ModuleError};
use opf_models::event::{send_event, Event};
use opf_models::{
    metadata::{Arg, Args},
    Target, TargetType,
};
use std::collections::HashMap;
use std::net::{SocketAddr, ToSocketAddrs};
use std::time::Duration;
use tokio::net::TcpStream;
use tokio::sync::mpsc::UnboundedSender;

const MOST_COMMON_PORTS: &[(u16, &str)] = &[
    (80, "HTTP"),
    (23, "Telnet"),
    (443, "HTTPS"),
    (21, "FTP"),
    (22, "SSH"),
    (25, "SMTP"),
    (3389, "RDP"),
    (110, "POP3"),
    (445, "SMB"),
    (139, "NetBIOS"),
    (143, "IMAP"),
    (53, "DNS"),
    (135, "RPC"),
    (3306, "MySQL"),
    (8080, "HTTP-alt"),
    (1723, "PPTP"),
    (111, "RPCbind"),
    (995, "POP3S"),
    (993, "IMAPS"),
    (5900, "VNC"),
    (1025, "MS-RPC"),
    (587, "SMTP-MSA"),
    (8888, "HTTP-proxy"),
    (199, "SNMP"),
    (1720, "H.323"),
    (465, "SMTPS"),
    (548, "AFP"),
    (113, "IDENT"),
    (81, "HTTP-alt"),
    (6001, "X11"),
    (10000, "Webmin"),
    (514, "Syslog"),
    (5060, "SIP"),
    (179, "BGP"),
    (1026, "Windows-RPC"),
    (2000, "Cisco-Sccp"),
    (8443, "HTTPS-alt"),
    (8000, "iRDMI"),
    (32768, "Filenet-TMS"),
    (554, "RTSP"),
    (26, "SMTP-alt"),
    (1433, "MS-SQL"),
    (49152, "Reserved"),
    (2001, "Cisco-SCCP"),
    (515, "LPD"),
    (8008, "HTTP-alt"),
    (49154, "Reserved"),
    (1027, "Windows-RPC"),
    (5666, "NRPE"),
    (646, "LDP"),
    (5000, "UPnP"),
    (5631, "PCAnywhere"),
    (631, "IPP"),
    (49153, "Reserved"),
    (8081, "HTTP-alt"),
    (2049, "NFS"),
    (88, "Kerberos"),
    (79, "Finger"),
    (5800, "VNC-HTTP"),
    (106, "POP3PW"),
    (2121, "FTP-proxy"),
    (1110, "POP3-alt"),
    (49155, "Reserved"),
    (6000, "X11"),
    (513, "Rlogin"),
    (990, "FTPS"),
    (5357, "WS-Discovery"),
    (427, "SLP"),
    (49156, "Reserved"),
    (543, "Klogin"),
    (544, "Kshell"),
    (5101, "TFTP"),
    (144, "NeWS"),
    (7, "Echo"),
    (389, "LDAP"),
    (8009, "AJP13"),
    (3128, "Squid-proxy"),
    (444, "SNPP"),
    (9999, "Abyss-Web"),
    (5009, "WinFS"),
    (7070, "RealServer"),
    (5190, "AIM"),
    (3000, "Cloud9-IDE"),
    (5432, "PostgreSQL"),
    (1900, "UPnP"),
    (3986, "MAPPER"),
    (13, "Daytime"),
    (1029, "Windows-RPC"),
    (9, "Discard"),
    (5051, "ITA"),
    (6646, "McAfee-Update"),
    (49157, "Reserved"),
    (1028, "Windows-RPC"),
    (873, "Rsync"),
    (1755, "MMS"),
    (2717, "PN-RTM"),
    (4899, "RAdmin"),
    (9100, "JetDirect"),
    (119, "NNTP"),
    (37, "Time"),
    (1000, "Cadlock"),
    (3001, "Nessus"),
    (5001, "SIP"),
    (82, "XFER"),
    (10010, "Zebra"),
    (1030, "BMC"),
    (9090, "WebSM"),
    (2107, "MSMQ"),
    (1024, "Reserved"),
    (2103, "MSMQ"),
    (6004, "X11"),
    (1801, "MSMQ"),
    (5050, "MMCC"),
    (19, "Chargen"),
    (8031, "ProEd"),
    (1041, "DanwareNetOp"),
    (255, "RIP"),
    (2967, "Symantec-AV"),
    (1049, "DanwareNetOp"),
    (1048, "Neptune"),
    (1053, "RemoteAssistant"),
    (3703, "Adobe-Server-3"),
    (1056, "VFO"),
    (1065, "Sybase"),
    (1064, "JSTO"),
    (1054, "BRVRead"),
    (17, "Quote"),
    (808, "CCProxy"),
    (3689, "DAAP"),
    (1031, "BBN-IAH"),
    (1044, "DCUtility"),
    (1071, "BSQUARE-VOIP"),
    (5901, "VNC-1"),
    (9102, "Bacula"),
    (100, "Newacct"),
    (8010, "XMPP"),
    (2869, "ICS"),
    (1039, "Streamlined"),
    (5120, "ISS-Agent"),
    (4001, "Cisco-Mgmt"),
    (9000, "CSlistener"),
    (2105, "MiniPay"),
    (636, "LDAPS"),
    (1038, "BBN-IAH"),
    (2601, "Zebra-IP"),
    (7000, "Afs3-Fileserver"),
    (1, "TCPMUX"),
    (1066, "FPO-FNS"),
    (1069, "Cognex-Insight"),
    (625, "AppleShare"),
    (311, "AppleShare"),
    (280, "HTTP-Mgmt"),
    (254, "SGMP"),
    (4000, "ICQ"),
    (5003, "FileMaker"),
    (1761, "cft"),
    (2002, "Globe"),
    (2005, "Orbix"),
    (1998, "X25-SVC"),
    (1032, "BBN-IAH"),
    (1050, "CORBA"),
    (6112, "BattleNet"),
    (3690, "SVN"),
    (1521, "Oracle"),
    (2161, "APC-Agent"),
    (6002, "X11"),
    (1080, "SOCKS"),
    (2401, "CVS"),
    (4045, "Solaris-LPD"),
    (902, "VMware-Auth"),
    (7937, "NSR-Listener"),
    (787, "QPASA-Agent"),
    (1058, "NIM"),
    (2383, "SQL-Analyzer"),
    (32771, "Filenet-RMI"),
    (1033, "BBN-IAH"),
    (1040, "Netarx"),
    (1059, "Forge-IA"),
    (50000, "SAP-Dispatcher"),
    (5555, "Personal-Agent"),
    (10001, "SCP-Config"),
    (1494, "Citrix-ICA"),
    (2301, "Compaq-HTTPS"),
    (593, "HTTP-RPC"),
    (3, "Compressnet"),
    (3268, "MS-Global-Catalog"),
    (7938, "Lanner-lic"),
    (1234, "VLC"),
    (1022, "RFC-Compliant"),
    (1035, "MX-Is"),
    (9001, "T2-Bravo"),
    (1074, "Jam-Link"),
    (8002, "Teradata"),
    (1036, "Nebula"),
    (1037, "AMS"),
    (464, "KPassword"),
    (1935, "RTMP"),
    (6666, "IRCD"),
    (2003, "Finger"),
    (497, "Retrospect"),
    (5601, "Kibana"),
    (9200, "Elasticsearch"),
    (9300, "Elasticsearch"),
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

        for (port, service) in MOST_COMMON_PORTS {
            let port = port.clone();
            let _ = tokio::spawn(PortScanner::scan_port(
                tx.clone(),
                target.clone(),
                port,
                service,
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

    async fn scan_port(tx: UnboundedSender<Event>, hostname: String, port: u16, service: &str, parent: String) {
        let timeout = Duration::from_secs(3);
        let socket_addresses: Vec<SocketAddr> =
            match format!("{}:{}", hostname, port).to_socket_addrs() {
                Ok(socket) => socket.collect(),
                Err(_) => return,
            };

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
                target.insert(String::from("service"), service.to_string());

                if let Ok(target) = Target::try_from(target) {
                    let _ = send_event(&tx, Event::ResultsModule(vec![target])).await;
                    let _ = send_event(
                        &tx,
                        Event::ResponseSimple(format!(
                            "hostname '{}' have a new port '{}' seem to be '{}'",
                            hostname, port, service
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
