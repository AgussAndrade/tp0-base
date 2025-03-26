import csv
import datetime
import time

""" Bets storage location. """
STORAGE_FILEPATH = "./bets.csv"
""" Simulated winner number in the lottery contest. """
LOTTERY_WINNER_NUMBER = 7574


""" A lottery bet registry. """
class Bet:
    def __init__(self, agency: str, first_name: str, last_name: str, document: str, birthdate: str, number: str):
        """
        agency must be passed with integer format.
        birthdate must be passed with format: 'YYYY-MM-DD'.
        number must be passed with integer format.
        """
        self.agency = int(agency)
        self.first_name = first_name
        self.last_name = last_name
        self.document = document
        self.birthdate = datetime.date.fromisoformat(birthdate)
        self.number = int(number)

""" Checks whether a bet won the prize or not. """
def has_won(bet: Bet) -> bool:
    return bet.number == LOTTERY_WINNER_NUMBER

"""
Persist the information of each bet in the STORAGE_FILEPATH file.
Not thread-safe/process-safe.
"""
def store_bets(bets: list[Bet]) -> None:
    with open(STORAGE_FILEPATH, 'a+') as file:
        writer = csv.writer(file, quoting=csv.QUOTE_MINIMAL)
        for bet in bets:
            writer.writerow([bet.agency, bet.first_name, bet.last_name,
                             bet.document, bet.birthdate, bet.number])

"""
Loads the information all the bets in the STORAGE_FILEPATH file.
Not thread-safe/process-safe.
"""
def load_bets() -> list[Bet]:
    with open(STORAGE_FILEPATH, 'r') as file:
        reader = csv.reader(file, quoting=csv.QUOTE_MINIMAL)
        for row in reader:
            yield Bet(row[0], row[1], row[2], row[3], row[4], row[5])


def recv_until(sock, delimiter=b'\t'):
    data = b''

    while True:
        chunk = sock.recv(1024)
        if not chunk:
            break
        data += chunk
        if delimiter in data:
            break
    return data

def construct_bet_by_msg(msg):
    return Bet(msg[0],msg[1],msg[2],msg[3],msg[4],msg[5])

def check_bet_msg(msg):
    if len(msg) != 6:
        return False
    if not msg[0].isdigit() or not msg[3].isdigit() or not msg[5].isdigit():
        return False
    if not msg[1] or not msg[2]:
        return False
    if not is_strict_iso_date_format(msg[4]):
        return False
    return True    

def is_strict_iso_date_format(date_str: str) -> bool:
    parts = date_str.split('-')
    if len(parts) != 3:
        return False

    year, month, day = parts

    if not (year.isdigit() and len(year) == 4):
        return False
    if not (month.isdigit() and len(month) == 2):
        return False
    if not (day.isdigit() and len(day) == 2):
        return False

    return True

def take_winners():
        all_bets = load_bets()
        winners_by_agency = {}

        for bet in all_bets:
            if has_won(bet):
                agency = bet.agency
                if agency not in winners_by_agency:
                    winners_by_agency[agency] = []
                winners_by_agency[agency].append(bet.document)
        
        return winners_by_agency

def format_winners_list(document_list):
    return ';'.join(document_list) + '\t'