use nom::bytes::complete::{tag, take_until};
use nom::combinator::map_res;
use nom::multi::fold_many0;
use nom::IResult;
use std::str::FromStr;
use strum_macros::{Display, EnumString};

#[derive(Debug, Default, Clone)]
pub struct Metadata {
    pub name: Option<String>,
    pub description: Option<String>,
    pub author: Option<String>,
    pub args: Vec<Arg>,
    pub extends: Vec<Requirements>,
}

#[derive(Debug, Clone)]
pub struct Arg {
    pub is_target: bool,
    pub is_optional: bool,
    pub name: String,
    pub value: Option<String>,
}

#[derive(Debug, PartialEq, EnumString, Display, Clone)]
#[strum(serialize_all = "lowercase")]
pub enum Field {
    Name,
    Description,
    Author,
    Args,
    Require,
}

#[derive(Debug, PartialEq, EnumString, Display, Clone)]
#[strum(serialize_all = "lowercase")]
pub enum Requirements {
    Http,
    Scraper,
    Target,
    Network,
    Common,
}

fn field(input: &str) -> IResult<&str, (Field, &str)> {
    let (input, _) = tag("-- ")(input)?;
    let (input, name) = map_res(take_until(": "), Field::from_str)(input)?;
    let (input, _) = tag(": ")(input)?;
    let (input, value) = take_until("\n")(input)?;
    let (input, _) = tag("\n")(input)?;

    Ok((input, (name, value)))
}

fn parse_fields(input: &str) -> IResult<&str, Vec<(Field, &str)>> {
    let (input, lines) = fold_many0(field, Vec::new, |mut acc: Vec<_>, item| {
        acc.push(item);
        acc
    })(input)?;
    let (input, _) = tag("\n")(input)?;

    Ok((input, lines))
}

pub fn parse(input: &str) -> Result<Metadata, String> {
    let mut metadata = Metadata::default();
    let results: Vec<(Field, &str)> = match parse_fields(input) {
        Ok(res) => res.1,
        Err(_) => return Err("can't parse metadata".to_string()),
    };

    for (meta, value) in results {
        match meta {
            Field::Name => metadata.name = Some(value.to_string()),
            Field::Description => metadata.description = Some(value.to_string()),
            Field::Author => metadata.author = Some(value.to_string()),
            Field::Args => {
                let elements = value.split(",").collect::<Vec<&str>>();
                for element in elements {
                    let mut name = element.trim();
                    let mut is_target = false;
                    let mut is_optional = false;
                    if name.contains("target") {
                        let start = match name.find(":") {
                            Some(size) => size + 1,
                            None => 0,
                        };
                        name = &name[start..];
                        is_target = true;
                    }

                    if name.contains("opt") {
                        let start = match name.find(":") {
                            Some(size) => size + 1,
                            None => 0,
                        };
                        name = &name[start..];
                        is_optional = true;
                    }

                    metadata.args.push(Arg {
                        is_target,
                        is_optional,
                        name: name.to_string(),
                        value: None,
                    });
                }
            }
            Field::Require => {
                let elements = value.split(",").collect::<Vec<&str>>();
                for element in elements {
                    let name = element.trim();
                    match Requirements::from_str(name) {
                        Ok(extend) => metadata.extends.push(extend),
                        Err(_) => return Err("requirements not available".to_string()),
                    }
                }
            }
        }
    }
    Ok(metadata)
}
