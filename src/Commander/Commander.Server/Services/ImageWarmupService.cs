using Commander.Infrastructure.Adapters;

namespace Commander.Server.Services;

public class ImageWarmupService(DockerRunnerAdapter adapter, ILogger<ImageWarmupService> logger)
    : IHostedService
{
    public async Task StartAsync(CancellationToken cancellationToken)
    {
        try
        {
            logger.LogInformation("Ensuring runner image is available...");
            await adapter.EnsureImageAsync();
            logger.LogInformation("Runner image ready.");
        }
        catch (Exception e)
        {
            logger.LogWarning(e, "Could not pre-pull Image. Next try on first job execution!");
        }
    }

    public Task StopAsync(CancellationToken cancellationToken) => Task.CompletedTask;
}
