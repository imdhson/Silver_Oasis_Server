from bardapi import Bard
from dotenv import load_dotenv
import os

load_dotenv(verbose=True)
token = os.getenv("BARD_API_KEY")

with open("CallBard/input.txt", "r") as f:
    input_text = f.read()
bard = Bard(token = token)
response = bard.get_answer(input_text)['content']

with open("CallBard/output.txt", "w") as f:
    f.writelines(response)