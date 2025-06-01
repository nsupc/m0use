import src.eurocore as eurocore
import src.ns as ns

from src.config import Config


def main():
    config = Config.from_yml()

    nations = ns.get_nations(
        config.user, config.region, config.cache_results, config.request_rate
    )

    eurocore.send_telegrams(
        nations,
        config.eurocore.url,
        config.eurocore.user,
        config.eurocore.password,
        config.telegram.id,
        config.telegram.key,
    )


if __name__ == "__main__":
    main()
