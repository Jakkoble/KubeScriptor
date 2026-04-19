pub mod orchestrator {
    tonic::include_proto!("orchestrator");
}

use std::error::Error;

use orchestrator::{SubmitJobRequest, orchestrator_service_client::OrchestratorServiceClient};
use tonic::transport::Channel;

pub struct CommanderClient {
    client: OrchestratorServiceClient<Channel>,
}

impl CommanderClient {
    pub async fn connect(addr: &str) -> Result<Self, Box<dyn Error>> {
        let client = OrchestratorServiceClient::connect(addr.to_string()).await?;

        Ok(Self { client })
    }

    pub async fn submit_job(&mut self, yaml_payload: String) -> Result<String, Box<dyn Error>> {
        let response = self
            .client
            .submit_job(SubmitJobRequest { yaml_payload })
            .await?;

        Ok(response.into_inner().job_id)
    }
}
