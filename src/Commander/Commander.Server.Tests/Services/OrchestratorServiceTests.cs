using Commander.Core.Entities;
using Commander.Core.Factories;
using Commander.Core.Ports;
using Grpc.Core;
using Microsoft.Extensions.Logging.Abstractions;
using Moq;

namespace Commander.Server.Tests.Services;

public class OrchestratorServiceTests
{
  private readonly JobDefinitionFactory _factory;
  private readonly Mock<IRunnerPort> _runnerPortMock;
  private readonly Server.Services.OrchestratorService _service;

  public OrchestratorServiceTests()
  {
    _factory = new();
    _runnerPortMock = new();

    var logger = NullLogger<Server.Services.OrchestratorService>.Instance;
    _service = new(_factory, _runnerPortMock.Object, logger);
  }

  [Fact]
  public async Task SubmitJob_WithValidYaml()
  {
    var validYaml = """
    name: my-super-test-job
    commands:
     - echo "Hello World"
     - ls -la
    """;

    var request = new SubmitJobRequest { YamlPayload = validYaml };
    var response = await _service.SubmitJob(request, null!);

    Assert.False(string.IsNullOrWhiteSpace(response.JobId));
    _runnerPortMock.Verify(
      port => port.ExecuteJob(It.Is<Job>(j =>
        j.Name == "my-super-test-job" &&
        j.Status == JobStatus.Running)),
      Times.Once);
  }

  [Fact]
  public async Task SubmitJob_WithoutJob()
  {
    var request = new SubmitJobRequest();
    var exception = await Assert.ThrowsAsync<RpcException>(() => _service.SubmitJob(request, null!));

    Assert.Equal(StatusCode.InvalidArgument, exception.StatusCode);
    _runnerPortMock.Verify(port => port.ExecuteJob(It.IsAny<Job>()), Times.Never);

  }

  [Fact]
  public async Task SubmitJob_WithoutName()
  {
    var incompleteYaml = """
    commands:
     - echo "Hello World"
    """;

    var request = new SubmitJobRequest { YamlPayload = incompleteYaml };
    var exception = await Assert.ThrowsAsync<RpcException>(() => _service.SubmitJob(request, null!));

    Assert.Equal(StatusCode.InvalidArgument, exception.StatusCode);
    _runnerPortMock.Verify(port => port.ExecuteJob(It.IsAny<Job>()), Times.Never);
  }
}
