use std::collections::HashMap;

pub fn get_params<'a>(
    params: &'a HashMap<String, String>,
    want: Vec<String>,
) -> Result<Vec<String>, String> {
    let mut ret: Vec<String> = vec![];
    for need in want {
        if let Some(n) = params.get(&need) {
            ret.push(n.clone());
            continue;
        }

        return Err(format!("Parameter {} not found in params lists", need));
    }

    Ok(ret)
}

pub fn get_param<'a>(params: &'a HashMap<String, String>, want: &'a str) -> Result<String, String> {
    if let Some(n) = params.get(want) {
        return Ok(n.clone());
    }

    return Err(format!("Parameter {} not found in params lists", want));
}
