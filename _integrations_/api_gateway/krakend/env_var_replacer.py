from shutil import move, copymode
from os import fdopen
import re, os

regex = "\{\s(\$.*?)\s\}"

with open('./krakend-final.json','w+') as new_file:
    with open('./krakend.json') as old_file:
        for line in old_file:
            exists = re.search(regex, line)
            if not exists:
                new_file.write(line)
                continue

            pattern = exists.group()
            env_var_name = pattern[3:len(pattern)-2]
            env_var_value = os.getenv(env_var_name)

            if env_var_value is None:
                raise Exception(f"the {env_var_name} secret is not set")

            new_file.write(line.replace(pattern, env_var_value))

    copymode('krakend.json', 'krakend-final.json')