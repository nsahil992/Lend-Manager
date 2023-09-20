friends = ["Courage", "Doreamon", "Ben"]         #friends array
lent = [[], ["Computer, Gadget"], ["Watch"]]     #lent array
while true
    println("What do you want to do?(i.e takeback/give/newfriend)")
    user_action = readline() #takes user action
end

if user_action == "takeback"
    elseif user_action == "give"
    elseif user_action == "newfriend"
else
    println("Sorry, I didn't understand that. Valid actions: takeback/give/newfriend")
end
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

