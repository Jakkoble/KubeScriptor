use crate::client::CommanderClient;

mod client;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let path = std::env::args()
        .nth(1)
        .unwrap_or_else(|| "job.yaml".to_string());

    let yaml_payload =
        std::fs::read_to_string(&path).unwrap_or_else(|_| panic!("File not found {path}"));

    println!("Sent job: {path}");

    let mut client = CommanderClient::connect("http://[::1]:5271").await?;

    let job_id = client.submit_job(yaml_payload).await?;

    println!("Job accepted. ID: {job_id}");

    Ok(())
}
