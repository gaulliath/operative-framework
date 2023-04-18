use std::collections::HashMap;
use std::ops::Add;

use async_trait::async_trait;

use reqwest::header;
use tokio::sync::mpsc::UnboundedSender;

use crate::CompiledModule;
use opf_models::error::{ErrorKind, Module as ModuleError};
use opf_models::event::Event;
use opf_models::{
    metadata::{Arg, Args},
    Target, TargetType,
};
use serde::{Deserialize, Serialize};
use serde_json::Value;

#[derive(Clone)]
pub struct InseeSearch {}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Root {
    #[serde(skip)]
    pub header: Header,
    pub etablissements: Vec<Etablissement>,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Header {
    pub statut: i64,
    pub message: String,
    pub total: i64,
    pub debut: i64,
    pub nombre: i64,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Etablissement {
    #[serde(skip)]
    pub siren: String,
    #[serde(skip)]
    pub nic: String,
    #[serde(skip)]
    pub siret: String,
    #[serde(skip)]
    pub statut_diffusion_etablissement: String,
    #[serde(skip)]
    pub date_creation_etablissement: String,
    #[serde(skip)]
    pub tranche_effectifs_etablissement: Value,
    #[serde(skip)]
    pub annee_effectifs_etablissement: Value,
    #[serde(skip)]
    pub activite_principale_registre_metiers_etablissement: Value,
    #[serde(skip)]
    pub date_dernier_traitement_etablissement: String,
    #[serde(skip)]
    pub etablissement_siege: bool,
    #[serde(skip)]
    pub nombre_periodes_etablissement: i64,
    #[serde(skip)]
    pub unite_legale: UniteLegale,
    pub adresse_etablissement: AdresseEtablissement,
    #[serde(rename = "adresse2Etablissement")]
    #[serde(skip)]
    pub adresse2etablissement: Adresse2Etablissement,
    #[serde(skip)]
    pub periodes_etablissement: Vec<PeriodesEtablissement>,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct UniteLegale {
    pub etat_administratif_unite_legale: String,
    pub statut_diffusion_unite_legale: String,
    pub date_creation_unite_legale: String,
    pub categorie_juridique_unite_legale: String,
    pub denomination_unite_legale: Value,
    pub sigle_unite_legale: Value,
    #[serde(rename = "denominationUsuelle1UniteLegale")]
    pub denomination_usuelle1unite_legale: Value,
    #[serde(rename = "denominationUsuelle2UniteLegale")]
    pub denomination_usuelle2unite_legale: Value,
    #[serde(rename = "denominationUsuelle3UniteLegale")]
    pub denomination_usuelle3unite_legale: Value,
    pub sexe_unite_legale: String,
    pub nom_unite_legale: String,
    pub nom_usage_unite_legale: Value,
    #[serde(rename = "prenom1UniteLegale")]
    pub prenom1unite_legale: String,
    #[serde(rename = "prenom2UniteLegale")]
    pub prenom2unite_legale: Value,
    #[serde(rename = "prenom3UniteLegale")]
    pub prenom3unite_legale: Value,
    #[serde(rename = "prenom4UniteLegale")]
    pub prenom4unite_legale: Value,
    pub prenom_usuel_unite_legale: String,
    pub pseudonyme_unite_legale: Value,
    pub activite_principale_unite_legale: String,
    pub nomenclature_activite_principale_unite_legale: String,
    pub identifiant_association_unite_legale: Value,
    pub economie_sociale_solidaire_unite_legale: Value,
    pub societe_mission_unite_legale: Value,
    pub caractere_employeur_unite_legale: String,
    pub tranche_effectifs_unite_legale: Value,
    pub annee_effectifs_unite_legale: Value,
    pub nic_siege_unite_legale: String,
    pub date_dernier_traitement_unite_legale: String,
    pub categorie_entreprise: String,
    pub annee_categorie_entreprise: String,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct AdresseEtablissement {
    pub complement_adresse_etablissement: Option<Value>,
    pub numero_voie_etablissement: Option<String>,
    pub indice_repetition_etablissement: Option<Value>,
    pub type_voie_etablissement: Option<String>,
    pub libelle_voie_etablissement: Option<String>,
    pub code_postal_etablissement: Option<String>,
    pub libelle_commune_etablissement: Option<String>,
    pub libelle_commune_etranger_etablissement: Option<Value>,
    pub distribution_speciale_etablissement: Option<Value>,
    pub code_commune_etablissement: Option<String>,
    pub code_cedex_etablissement: Option<Value>,
    pub libelle_cedex_etablissement: Option<Value>,
    pub code_pays_etranger_etablissement: Option<Value>,
    pub libelle_pays_etranger_etablissement: Option<Value>,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Adresse2Etablissement {
    #[serde(rename = "complementAdresse2Etablissement")]
    pub complement_adresse2etablissement: Value,
    #[serde(rename = "numeroVoie2Etablissement")]
    pub numero_voie2etablissement: Value,
    #[serde(rename = "indiceRepetition2Etablissement")]
    pub indice_repetition2etablissement: Value,
    #[serde(rename = "typeVoie2Etablissement")]
    pub type_voie2etablissement: Value,
    #[serde(rename = "libelleVoie2Etablissement")]
    pub libelle_voie2etablissement: Value,
    #[serde(rename = "codePostal2Etablissement")]
    pub code_postal2etablissement: Value,
    #[serde(rename = "libelleCommune2Etablissement")]
    pub libelle_commune2etablissement: Value,
    #[serde(rename = "libelleCommuneEtranger2Etablissement")]
    pub libelle_commune_etranger2etablissement: Value,
    #[serde(rename = "distributionSpeciale2Etablissement")]
    pub distribution_speciale2etablissement: Value,
    #[serde(rename = "codeCommune2Etablissement")]
    pub code_commune2etablissement: Value,
    #[serde(rename = "codeCedex2Etablissement")]
    pub code_cedex2etablissement: Value,
    #[serde(rename = "libelleCedex2Etablissement")]
    pub libelle_cedex2etablissement: Value,
    #[serde(rename = "codePaysEtranger2Etablissement")]
    pub code_pays_etranger2etablissement: Value,
    #[serde(rename = "libellePaysEtranger2Etablissement")]
    pub libelle_pays_etranger2etablissement: Value,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct PeriodesEtablissement {
    pub date_fin: Option<String>,
    pub date_debut: String,
    pub etat_administratif_etablissement: String,
    pub changement_etat_administratif_etablissement: bool,
    #[serde(rename = "enseigne1Etablissement")]
    pub enseigne1etablissement: Value,
    #[serde(rename = "enseigne2Etablissement")]
    pub enseigne2etablissement: Value,
    #[serde(rename = "enseigne3Etablissement")]
    pub enseigne3etablissement: Value,
    pub changement_enseigne_etablissement: bool,
    pub denomination_usuelle_etablissement: Value,
    pub changement_denomination_usuelle_etablissement: bool,
    pub activite_principale_etablissement: String,
    pub nomenclature_activite_principale_etablissement: String,
    pub changement_activite_principale_etablissement: bool,
    pub caractere_employeur_etablissement: String,
    pub changement_caractere_employeur_etablissement: bool,
}

#[async_trait]
impl CompiledModule for InseeSearch {
    fn name(&self) -> String {
        "insee.search".to_string()
    }

    fn author(&self) -> String {
        "Tristan Granier".to_string()
    }

    fn resume(&self) -> String {
        "Extract information from SIREN (FR)".to_string()
    }

    fn args(&self) -> Vec<Arg> {
        vec![
            Arg::new("target_id", true, false, None),
            Arg::new("target", false, false, None),
            Arg::new("siren", false, false, None),
        ]
    }

    fn target_type(&self) -> TargetType {
        TargetType::Company
    }

    async fn run(
        &self,
        params: Args,
        _tx: Option<UnboundedSender<Event>>,
    ) -> Result<Vec<Target>, ErrorKind> {
        let target_id = params
            .get("target_id")
            .ok_or(ErrorKind::Module(ModuleError::TargetNotAvailable))?
            .value
            .ok_or(ErrorKind::Module(ModuleError::TargetNotAvailable))?;
        let siren = params
            .get("siren")
            .ok_or(ErrorKind::Module(ModuleError::ParamNotAvailable(
                "siren".to_string(),
            )))?
            .value
            .ok_or(ErrorKind::Module(ModuleError::ParamNotAvailable(
                "siren".to_string(),
            )))?;

        println!("siren = {}", siren);
        let mut headers = header::HeaderMap::new();
        headers.insert(
            "User-Agent",
            "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/110.0"
                .parse()
                .unwrap(),
        );
        headers.insert("Accept", "application/json".parse().unwrap());
        headers.insert("Accept-Language", "en-US,en;q=0.5".parse().unwrap());
        headers.insert("Content-Type", "application/json".parse().unwrap());
        headers.insert(
            "Authorization",
            "Bearer 690e00d0-bd03-363d-8035-140e426b191e"
                .parse()
                .unwrap(),
        );
        headers.insert("Connection", "keep-alive".parse().unwrap());

        let client = reqwest::Client::builder()
            .redirect(reqwest::redirect::Policy::none())
            .build()
            .map_err(|e| ErrorKind::Module(ModuleError::Execution(e.to_string())))?;

        let url = format!(
            "https://api.insee.fr/entreprises/sirene/V3/siret?q=siren:{}&nombre=999",
            siren
        );
        println!("url => {}", url);
        let res = client
            .get(url)
            .headers(headers)
            .send()
            .await
            .map_err(|e| ErrorKind::Module(ModuleError::Execution(e.to_string())))?;

        let res: Root = res
            .json()
            .await
            .map_err(|e| ErrorKind::Module(ModuleError::Execution(e.to_string())))?;

        let mut results = vec![];
        for element in res.etablissements {
            let address = element.adresse_etablissement;
            let final_address = address
                .numero_voie_etablissement
                .unwrap_or("".to_string())
                .clone()
                .add(" ")
                .add(
                    address
                        .type_voie_etablissement
                        .unwrap_or("".to_string())
                        .as_str(),
                )
                .add(" ")
                .add(
                    address
                        .libelle_voie_etablissement
                        .unwrap_or("".to_string())
                        .as_str(),
                )
                .add(" ")
                .add(
                    address
                        .code_postal_etablissement
                        .unwrap_or("".to_string())
                        .as_str(),
                )
                .add(" ")
                .add(
                    address
                        .libelle_commune_etablissement
                        .unwrap_or("".to_string())
                        .as_str(),
                );

            let mut result = HashMap::new();

            result.insert(String::from("name"), final_address);
            result.insert(String::from("type"), String::from("address"));
            result.insert(String::from("parent"), target_id.clone());

            results.push(Target::try_from(result)?);
        }

        Ok(results)
    }
}

impl InseeSearch {
    pub fn new() -> Box<dyn CompiledModule> {
        Box::new(InseeSearch {})
    }
}
