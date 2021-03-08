# Golang-Challenge
Challenge test

We ask that you complete the following challenge to evaluate your development skills.

## The Challenge
Finish the implementation of the provided Transparent Cache package.

## Show your work

1.  Create a **Private** repository and share it with the recruiter ( please dont make a pull request, clone the private repository and create a new private one on your profile)
2.  Commit each step of your process so we can follow your thought process.
3.  Give your interviewer access to the private repo

## What to build
Take a look at the current TransparentCache implementation.

You'll see some "TODO" items in the project for features that are still missing.

The solution can be implemented either in Golang or Java ( but you must be able to read code in Golang to realize the exercise ) 

Also, you'll see that some of the provided tests are failing because of that.

The following is expected for solving the challenge:
* Design and implement the missing features in the cache
* Make the failing tests pass, trying to make none (or minimal) changes to them
* Add more tests if needed to show that your implementation really works
 
## Deliverables we expect:
* Your code in a private Github repo
* README file with the decisions taken and important notes

## Time Spent
We suggest not to spend more than 2 hours total, which can be done over the course of 2 days.  Please make commits as often as possible so we can see the time you spent and please do not make one commit.  We will evaluate the code and time spent.
 
What we want to see is how well you handle yourself given the time you spend on the problem, how you think, and how you prioritize when time is insufficient to solve everything.

Please email your solution as soon as you have completed the challenge or the time is up.


-----------------------------

## Decisions taken and important notes

### Decisions
First task was to quick read some general information about a "Transparent cache"

Second task was to read the repo to understand what it does, and where are the most important things

Once that I read the repo I identified 2 main task to fix, the maxAge and the parallel calls.
My decision was to start with the maxAge, because is the main functionality of the cache, the parallel 
calls are optimization, and it not makes sense been fast but running in the wrong direction.
(Time: ~0:21hs)

Third task was to fix the maxAge, I answered my clues about which was the expected behavior reading the comment
of the function `GetPriceFor()` that said `// GetPriceFor gets the price for the item, either from the cache or the 
actual service if it was not cached or too old` and the test `TestGetPriceFor_DoesNotReturnOldResults()` that gives a
 good example about what happen in every moment.
 To solve the maxAge of every price I created a new struct with an amount and an updatedAt, that gives me the opportunity
 of record the last update time of every price item.
 
Ford task was to fix the tests of maxAge, and add the necessary behaviors to mock the results with the new price struct.
The respective PR of the task 3 and 4 is:
https://github.com/santiagomoranlabat/challenge-deviget/pull/1
 (Time: ~0:50hs)

Five task was to add the feature of parallel calls, to make this possible I implemented one more layer of struct called
`PriceModel`, this was necessary because channels can return just one parameter. So `PriceMolde` returns a `Price` and 
an `Error`. With my new struct layer y started to implement a call with wg and I changed the return parameters of 
`GetPriceFor()` to `PriceModel`

Six task was to fix the test to make them understand the `PriceModel` layer. The respective PR of task 5 and 6 is: 
https://github.com/santiagomoranlabat/challenge-deviget/pull/2 
 (Time: ~0:18hs)

Seven task was to create the documentation. This was after my job day, and it took (Time: ~0:09hs): 
https://github.com/santiagomoranlabat/challenge-deviget/pull/3

####Total Time ~1:28hs

## Next steps

- Refactor of the functions in MVC folders
- Test to test more race conditions
- Test to test fail cases
