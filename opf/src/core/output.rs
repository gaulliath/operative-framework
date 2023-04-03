use super::Session;
use comfy_table::Table;

impl Session {
    pub fn output_table(&self, headers: Vec<String>, rows: Vec<Vec<String>>) {
        let mut table = Table::new();
        table.set_header(headers);
        table.add_rows(rows);
        println!("{table}");
    }
}
