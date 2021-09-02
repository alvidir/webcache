import os
import logging
import threading
import constants
import browser
import server

from dotenv import load_dotenv

if __name__ == "__main__":
    load_dotenv()
    logging.basicConfig(format=constants.LOG_FORMAT, level=logging.INFO,
                        datefmt=constants.TIME_FORMAT)

    config_path = os.getenv(constants.ENV_CONFIG_PATH)
    sleep_time = os.getenv(constants.ENV_SLEEP_TIME)
    sleep_time = float(sleep_time) if sleep_time else None
    
    fn_args = (config_path, sleep_time)
    fn_config = threading.Thread(target=browser.run, args=fn_args, daemon=True)
    fn_config.start()

    server_port = os.getenv(constants.ENV_SERVER_PORT)
    server.run(server_port)