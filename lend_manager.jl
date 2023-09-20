friends = ["Courage", "Doreamon", "Ben"]         #friends array
lent = [[], ["Computer, Gadget"], ["Watch"]]     #lent array
while true
    println("What do you want to do?(i.e takeback/give/newfriend/quit)")
    user_action = readline() #takes user action
    if user_action == "quit"
        break
    end
    if user_action == "takeback"
        println("These are your friends")
        for friend in friends
            println(friend)
        end

        println("Which friend did you lent to? ")
        friend_name = readline()
        friend_index = findall(x -> x == friend_name, friends)
        if length(friend_index) == 0
            println("Sorry, I didn't find that friend")
            continue
        else
            friend_index = friend_index[1]
            if(length(lent[friend_index])) == 0
                println("You haven't given anything to $(friend_name)")
                continue
            end
            println("This is what you gave to $(friend_name)")
            for item in lent[friend_index]
                println(item)
            end
            println("What did you takeback from $(friend_name)")
            item_name = readline()
             item_index = findall(x -> x == item_name, lent[friend_index])
              if length(item_index) == 0
                println("Sorry, I did'nt find that item.")
                continue
            else
                item_index = item_index[1]
                 deleteat!(lent[friend_index], item_index)
                    println("Alright, I'll remember that you took $(item_name) from $(friend_name)")
            end
        end
        elseif  user_action == "give"
         println("These are your friends:")
           for friend in friends
            println(friend)
        end
        println("Which friend did you lend to? ")
         friend_name = readline()
         friend_index = findall(x -> x == friend_name, friends)
         if length(friend_index) == 0
            println("Sorry, I didn't find that friend.")
        continue
    else
        friend_index = friend_index[1]
        println("What did you lend to $(friend_name)?")
        item_name = readline()
        push!(lent[friend_index], item_name)
        println("Got it! You lent $(item_name) to $(friend_name).")
    end
elseif user_action == "newfriend"
    println("Who is your new friend?")
    friend_name = readline()
    push!(friends, friend_name)
    push!(lent, [])
else
    println("Sorry, I didn't understand that. (Valid choices: give/takeback/newfriend):")
end
end




