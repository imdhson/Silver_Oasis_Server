import bardapi
from dotenv import load_dotenv
import os
import json

load_dotenv(verbose=True)
token = os.getenv("BARD_API_KEY")

with open("CallBard/input.txt", "r") as f:
    input_text = f.read()

response = bardapi.core.Bard(token).get_answer(input_text)['content']

with open("CallBard/output.txt", "w") as f:
    f.writelines(response)