import matplotlib.pyplot as plt
import math


def median(args):
    args.sort()
    length = len(args)
    middle = length // 2
    if length & 1 == 1:
        return args[middle]
    return (args[middle - 1] + args[middle]) // 2


def sequence(n):
    # return sum(i for i in range(1, n+1))
    return n*(n+1)//2


def cost(lower, upper, data) -> list:
    return [sum(sequence(abs(j - i)) for j in data) for i in range(lower, upper+1)]


def main():
    data = [int(i) for i in open("example.txt", "r").read().split(",")]
    print(data)
    bounds = min(data), max(data)
    fuel = cost(*bounds, data)
    print(fuel)
    plt.plot(range(bounds[0], bounds[1] + 1), fuel)
    plt.scatter(range(bounds[0], bounds[1] + 1), fuel)
    plt.show()


if __name__ == "__main__":
    main()
