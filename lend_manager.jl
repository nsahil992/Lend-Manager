lent = Dict{String, Array{String}}("Courage" => [], "Ben" => ["Watch"], "Mr.Bean" => ["Teddy"])
while true
    println("What do you want to do?(i.e takeback/give/newfriend/quit)")
    user_action = readline() #takes user action
    if user_action == "quit"
        break
    end

    #if user inputs quit then break the loop and end it

    if user_action == "takeback"
        println("These are your friends")
        for friend in keys(lent)
            println(friend)
        end

    # if user inputs takeback then he has given something to the friend
    #This part will show him all his friends

        println("Which friend did you lent to? ")
        friend_name = readline()
        friend_index = findall(x -> x == friend_name, friends)
        if !(friend_name in keys(lent))
            println("Sorry, I didn't find that friend")
            continue

    # if the user wants to takeback something, the compiler should know which friend did he lent to
    # so the compiler will take the name of the friend from user and will check whether the friend exists or not
    #if the index of the friend is 0 then the friend doesn't exist in the friends array

        else
            if(length(lent[friend_name])) == 0
                println("You haven't given anything to $(friend_name)")
                continue
            end

    # if we found the friend in the array, there might be a possiblity that we haven't lent anything to the friend
    # hence, in line 34, if we the item that we have lent doesn't exist then we haven't lent anything

            println("This is what you gave to $(friend_name)")
            for item in lent[friend_name]
                println(item)
            end

    # if we found the item that we have lent to the friend then it will display the friend name to whom we lent to and item

            println("What did you takeback from $(friend_name)")
            item_name = readline()
             item_index = findall(x -> x == item_name, lent[friend_name])
              if length(item_index) == 0
                println("Sorry, I did'nt find that item.")
                continue

    # if we have given multiple items to the same friend then it will take the name of the item from the user
    #if the if the name of the item from multiple item doesn't exist then it won't give you the items name
            else
                item_index = item_index[1]
                 deleteat!(lent[friend_name], item_index)
                    println("Alright, I'll remember that you took $(item_name) from $(friend_name)")
            end
        end
    
    #if it found the item then it will delete the item from the lent array as you are taking the item back

        elseif  user_action == "give"
         println("These are your friends:")
           for friend in keys(lent)
            println(friend)
        end

    #if the user inputs give then it will display the name of the user's friends

        println("Which friend do you want to lend to? ")
         friend_name = readline()
         friend_index = findall(x -> x == friend_name, friends)
         if !(friend_name in keys(lent))
            println("Sorry, I didn't find that friend.")
        continue

    #it will ask the user to which friend he has to lend took
    #it will take the name of the friend from the user and if the friend doesn't exist in the friend array, it won't give any name
    else
        friend_index = friend_index[1]
        println("What do you want to lend to $(friend_name)?")
        item_name = readline()
        push!(lent[friend_name], item_name)
        println("Got it! You lent $(item_name) to $(friend_name).")
    end

    #if it found the friend then it will take the item name from the user and push the item in the array

elseif user_action == "newfriend"
    println("Who is your new friend?")
    friend_name = readline()
    lent[friend_name] = []

    #if user inputs new friend then it will add the new friends name in the array
else
    println("Sorry, I didn't understand that. (Valid choices: give/takeback/newfriend):")

    #if the user inputs something that is not valid then it will ask user for valid options
end
end
println("bye. . .")


