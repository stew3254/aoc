from functools import reduce


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


# Convert a digit to the sorted list of segments
def convert_digit(digit, master: dict) -> str:
    return "".join(sorted(master[i] for i in digit))


def guess_digit(digit: set, count: dict, master: dict):
    # Guess the digit based on the length of the set
    match len(digit):
        case 2:
            # Solve digit 1
            for i, segment in enumerate(digit):
                # This is the bottom right segment
                if count[segment] == 9:
                    # Fix top right segment
                    master[list(digit)[i ^ 1]] = "c"
                    return 1
        case 3:
            # Solve the segment problem for digit 7
            for i, segment in enumerate(digit):
                v = master[segment]
                if type(v) is set:
                    # Fix top segment
                    master[segment] = "a"
            return 7
        case 4:
            # Get missing segment
            segment = list(digit - {k for k, v in master.items() if v in {"b", "c", "f"}})[0]
            # Fix middle segment
            master[segment] = "d"
            return 4
        case 5:
            # Could be 2, 3 or 5

            # Get missing segment
            segment = [k for k, v in master.items() if type(v) is set]
            if len(segment) > 0:
                master[segment[0]] = "g"

            # Convert digit to easier to check format
            d = convert_digit(digit, master)
            match d:
                case "acdeg":
                    return 2
                case "acdfg":
                    return 3
                case "abdfg":
                    return 5
        case 6:
            # Convert digit to easier to check format
            d = convert_digit(digit, master)
            match d:
                case "abcefg":
                    return 0
                case "abdefg":
                    return 6
                case "abcdfg":
                    return 9
        case 7:
            return 8


def main():
    inp = open("input.txt").readlines()
    total = 0
    for line in inp:
        master = {k: None for k in "abcdefg"}
        digits, display = parse_input(line)
        digits.sort(key=lambda x: len(x))
        segment_count = {}
        for digit in digits:
            for segment in digit:
                if segment_count.get(segment) is None:
                    segment_count[segment] = 1
                else:
                    segment_count[segment] += 1
        guess_segments(segment_count, master)
        # Solve digits
        [guess_digit(digit, segment_count, master) for digit in digits]
        total += reduce(lambda x, y: x * 10 + y, (guess_digit(digit, segment_count, master) for digit in display))
    #     total += sum(map(lambda x: 1, filter(lambda x: x is not None, (guess_digit(i, master) for i in display))))
    print(total)


if __name__ == "__main__":
    main()
