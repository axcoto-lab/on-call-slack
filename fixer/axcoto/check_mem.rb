class Oncall::Fixer::Axcoto::CheckMem

  def perform
    local do
      `ls`
    end

    ssh do
      
    end
  end
end
