def parse_input(line):
    parts = line.split("|")
    return [[set(i.strip()) for i in part.split()] for part in parts]


def guess_segments(count: dict, master: dict):
    for k, v in count.items():
        match v:
            case 4:
                master[k] = "e"
            case 6:
                master[k] = "b"
            case 7:
                master[k] = {"d", "g"}
            case 8:
                master[k] = {"a", "c"}
            case 9:
                master[k] = "f"


def guess_digit(digit: set, count: dict, master: dict):
    # Guess the digit based on the length of the set
    match len(digit):
        # Solve digit 1
        case 2:
            for i, segment in enumerate(digit):
                # This is the bottom right segment
                if count[segment] == 9:
                    # Get the other element in the digit and set it to c
                    master[list(digit)[i ^ 1]] = "c"
                    # Set other master key to its corresponding location
                    master[list(filter(lambda k: count[k] == 8 and type(master[k]) is set, count))[0]] = "a"
                    return 1
        case 3:
            return 7
        case 4:
            for i, segment in enumerate(digit):
                # This is the bottom right segment
                if count[segment] == 9:
                    # Get the other element in the digit and set it to c
                    master[list(digit)[i ^ 1]] = "c"
                    # Set other master key to its corresponding location
                    master[list(filter(lambda k: count[k] == 8 and type(master[k]) is set, count))[0]] = "a"
                    return 1
            return 4
        case 7:
            return 8


def main():
    inp = open("input.txt").readlines()
    total = 0
    for line in inp:
        master = {k: None for k in "abcdefg"}
        digits, display = parse_input(line)
        segment_count = {}
        for digit in digits:
            for segment in digit:
                if segment_count.get(segment) is None:
                    segment_count[segment] = 1
                else:
                    segment_count[segment] += 1
        guess_segments(segment_count, master)
        for digit in digits:
            guess_digit(digit, segment_count, master)
    #     total += sum(map(lambda x: 1, filter(lambda x: x is not None, (guess_digit(i, master) for i in display))))
    # print(total)


if __name__ == "__main__":
    main()
