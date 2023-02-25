import random

def populate_accounts(num_of_records):
    file = open("accounts.txt", "w")
    for num in range(num_of_records+1):
        query = f"INSERT INTO accounts(balance) VALUES ({round(random.uniform(0, 99999999.99), 2)});"
        file.write(query + '\n\n')
    file.close()

populate_accounts(100000)