import re
import random

def unknown():
    response = ["I didn't get you.",
                "I can't help you with that.",
                "Oops.. What does that mean?"][
        random.randrange(3)]
    return response



def message_probability(user_message, recognised_words, required_words=[]):
    message_certainty = 0
    has_required_words = True

    # Count matching words in message
    for word in user_message:
        if word in recognised_words:
            message_certainty += 1

    # Percent of recognised words
    percentage = float(message_certainty) / float(len(recognised_words))

    # Checks that the required words are in the string
    for word in required_words:
        if word not in user_message:
            has_required_words = False
            break
    # Return probability
    if has_required_words:
        return int(percentage * 100)
    else:
        return 0


def check_all_messages(message):
    #print (message)
    highest_prob_list = {}

    # Add response to dict
    def response(bot_response, list_of_words, required_words=[]):
        highest_prob_list[bot_response] = message_probability(message, list_of_words, required_words)

    # Predefined response
    response('Hello!', ['hello', 'hi', 'hey'])
    response('Bye! See you later', ['bye', 'goodbye'])
    response('I\'m doing fine.', ['how', 'are', 'you', 'doing'], required_words=['how','you'])
    response('You\'re welcome!', ['thank', 'thanks', 'thnx', 'thank you', 'thanku'])
    response('Thank You!', ['i', 'am', 'impressed'], required_words=['impressed'])

    best_match = max(highest_prob_list, key=highest_prob_list.get)
    return unknown() if highest_prob_list[best_match] < 1 else best_match


# Find response
def find_response(user_input):
    split_message = re.split(r'\s+|[,;?!.-]\s*', user_input.lower())
    response = check_all_messages(split_message)
    return response


# Running bot
while True:
    print('Bot: ' + find_response(input('User: ')))
