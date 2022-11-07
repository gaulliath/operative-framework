use dns_lookup;

/// checking if domain is valid
fn is_domain_valid<'a>(ctx: &'a mut rlua::Context) {

    let valid_domain = match ctx.create_function(|this, domain: String| {
        if let Ok(_) = addr::parse_domain_name(&domain) {
            return Ok(true);
        }
        Ok(false)
    }) {
        Ok(f) => f,
        Err(_) => {
            println!("function log disabled ?");
            return;
        }
    };

    ctx.globals().set("is_valid_domain", valid_domain).unwrap();
}

fn get_ip_address<'a>(ctx: &'a mut rlua::Context) {
    let get_ip_address = match ctx.create_function(|this, domain: String| {
        let ips = match dns_lookup::lookup_host(domain.as_str()) {
            Ok(ips) => ips,
            Err(_) => return Ok(vec![]),
        };

        let mut ip_addresses = vec![];
        for ip in ips {
            ip_addresses.push(ip.to_string());
        }

        Ok(ip_addresses)
    }) {
        Ok(f) => f,
        Err(_) => {
            println!("function log disabled ?");
            return;
        }
    };

    ctx.globals().set("get_ip", get_ip_address).unwrap();
}

pub fn extends_network<'a, 'b>(ctx: &'a mut rlua::Context) {
    // extends here
    is_domain_valid(ctx);
    get_ip_address(ctx);
}
